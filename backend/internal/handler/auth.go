package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
	"github.com/anvar-sharipov/telecom-map/internal/service"
	"github.com/anvar-sharipov/telecom-map/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo         *postgres.UserRepository
	RefreshTokenRepo *repository.RefreshTokenRepository
	AuthService      *service.AuthService
}

// UserResponse — структура для ответа (без пароля!)
type UserResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// toUserResponse конвертирует domain.User в UserResponse (убирает пароль)
func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

// ---------------- REGISTER ----------------
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	fmt.Println("GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")

	var req struct {
		Fullname        string `json:"fullname"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.NewBadRequest("Invalid request body")
	}

	if req.Password != req.ConfirmPassword {
		return utils.NewBadRequest("passwords do not match")
	}

	if req.Password == "" {
		return utils.NewBadRequest("password cant be empty")
	}

	if req.Username == "" {
		return utils.NewBadRequest("username cant be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewInternal("Failed to hash password")
	}

	newUser := &domain.User{
		FullName: req.Fullname,
		Username: req.Username,
		Password: string(hashed),
		IsActive: true,
	}

	if err := h.UserRepo.Create(newUser); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return utils.NewConflict("username already exists")
		}
		return utils.NewInternal("internal server error")
	}

	// ✅ Автологин через сервис
	tokens, err := h.AuthService.CreateSession(r.Context(), service.SessionData{
		UserID:    newUser.ID,
		UserAgent: r.UserAgent(),
		IP:        utils.GetClientIP(r),
	})
	if err != nil {
		return utils.NewInternal("failed to create session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true при HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  tokens.ExpiresAt,
	})

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message":      "user registered successfully",
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
		"user":         toUserResponse(newUser),
	})

	return nil
}

// ---------------- LOGIN ----------------
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.NewBadRequest("Invalid request body")
	}

	if req.Username == "" {
		return utils.NewBadRequest("username cant be empty")
	}

	if req.Password == "" {
		return utils.NewBadRequest("password cant be empty")
	}

	user, err := h.UserRepo.GetByUsername(req.Username)
	if err != nil {
		return utils.NewUnauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.NewUnauthorized("invalid credentials")
	}

	// ✅ Создаём сессию через сервис
	tokens, err := h.AuthService.CreateSession(r.Context(), service.SessionData{
		UserID:    user.ID,
		UserAgent: r.UserAgent(),
		IP:        utils.GetClientIP(r),
	})
	if err != nil {
		return utils.NewInternal("failed to create session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true при HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  tokens.ExpiresAt,
	})

	fmt.Println(toUserResponse(user))

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
		"user":         toUserResponse(user),
	})

	return nil
}

// ---------------- REFRESH ----------------
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return utils.NewUnauthorized("refresh token missing")
	}

	refreshToken := cookie.Value
	refreshTokenHash := utils.HashToken(refreshToken)

	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshTokenHash)
	if err != nil {
		return utils.NewUnauthorized("invalid refresh token")
	}

	currentAgent := r.UserAgent()
	if rt.UserAgent != currentAgent {
		return utils.NewUnauthorized("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = h.RefreshTokenRepo.Delete(r.Context(), refreshTokenHash)
		return utils.NewUnauthorized("refresh token expired")
	}

	// ✅ Ротация сессии через сервис
	tokens, err := h.AuthService.RotateSession(r.Context(), service.SessionData{
		UserID:    rt.UserID,
		UserAgent: r.UserAgent(),
		IP:        utils.GetClientIP(r),
	})
	if err != nil {
		return utils.NewInternal("failed to rotate session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  tokens.ExpiresAt,
	})

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
	})

	return nil
}

// ---------------- LOGOUT ----------------
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Println("No refresh cookie found")
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "logged out",
		})
		return nil
	}

	refreshHash := utils.HashToken(cookie.Value)
	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshHash)
	if err == nil {
		// ✅ Инвалидация сессии через сервис
		_ = h.AuthService.InvalidateSession(r.Context(), rt.UserID, r.UserAgent())
	} else {
		fmt.Println("Refresh token not found in DB, maybe already deleted")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})

	return nil
}

// ---------------- ME ----------------
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return utils.NewUnauthorized("not logged in")
	}

	refreshHash := utils.HashToken(cookie.Value)
	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshHash)
	if err != nil {
		return utils.NewUnauthorized("not logged in")
	}

	// 🔥 УЛУЧШЕНИЕ: возвращаем полные данные пользователя
	user, err := h.UserRepo.GetByID(rt.UserID)
	if err != nil {
		return utils.NewInternal("failed to get user")
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"user": toUserResponse(user), // 🔥 ПОЛНЫЕ ДАННЫЕ
	})

	// utils.WriteJSON(w, http.StatusOK, map[string]any{
	// 	"user_id":    rt.UserID,
	// 	"user_agent": rt.UserAgent,
	// })

	return nil
}
