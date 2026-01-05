package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
	"github.com/anvar-sharipov/telecom-map/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func getClientIP(r *http.Request) string {
	// 1. Если есть proxy
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// XFF может быть "client, proxy1, proxy2"
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	// 2. Без proxy
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

type AuthHandler struct {
	UserRepo         *postgres.UserRepository
	RefreshTokenRepo *repository.RefreshTokenRepository
}

// ---------------- REGISTER ----------------
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

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

	token, err := utils.GenerateToken(newUser.ID)
	if err != nil {
		return utils.NewInternal("Failed to generate token")
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
		"token":   token,
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

	// 1️⃣ ACCESS TOKEN
	accessToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		return utils.NewInternal("failed to generate access token")
	}

	// 2️⃣ REFRESH TOKEN
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return utils.NewInternal("failed to generate refresh token")
	}

	// 3️⃣ HASH refresh token
	refreshTokenHash := utils.HashToken(refreshToken)

	// 4️⃣ EXPIRE: срок действия, текущий момент времени + 7 дней
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// 5️⃣ SAVE refresh token
	userAgent := r.UserAgent()

	ip := getClientIP(r)

	err = h.RefreshTokenRepo.DeleteByUserAndAgent(
		r.Context(),
		user.ID,
		userAgent,
	)
	if err != nil {
		return utils.NewInternal("failed to cleanup old refresh token")
	}

	err = h.RefreshTokenRepo.Create(
		r.Context(),
		user.ID,
		refreshTokenHash,
		userAgent,
		ip,
		expiresAt,
	)
	if err != nil {
		log.Println("refresh token save error:", err)
		return utils.NewInternal("failed to save refresh token")
	}

	// ❌ БЫЛО
	// 6️⃣ RESPONSE
	// utils.WriteJSON(w, http.StatusOK, map[string]any{
	// 	"access_token":  accessToken,
	// 	// "refresh_token": refreshToken,
	// 	"expires_at":    expiresAt,
	// })

	// ✅ СТАНЕТ
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true при HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
	// SameSite — Strict ОК, но знай нюанс
	// ✔ Для SPA + same domain — отлично
	// 	⚠ Если потом будет:
	// 	фронт на другом домене
	// 	мобильное приложение
	// 	OAuth
	// 	👉 тогда меняют на:
	// 	SameSite: http.SameSiteLaxMode
	// 	Пока оставляем Strict, ты сделал правильно.

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
		"expires_at":   expiresAt,
	})
	return nil
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	// 🛠 ШАГ 2. REFRESH — читаем token из cookie
	// ❌ БЫЛО
	// ШАГ 3️⃣ — читаем refresh_token из body
	// var req struct {
	// 	RefreshToken string `json:"refresh_token"`
	// }

	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	return utils.NewBadRequest("invalid request body")
	// }

	// if req.RefreshToken == "" {
	// 	return utils.NewBadRequest("refresh token is required")
	// }

	// ✅ СТАНЕТ
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return utils.NewUnauthorized("refresh token missing")
	}

	// 📌 Body больше не нужен вообще
	refreshToken := cookie.Value

	// ШАГ 4️⃣ — хешируем refresh token
	refreshTokenHash := utils.HashToken(refreshToken)

	// ШАГ 5️⃣ — ищем refresh token в БД, ❗ В БД у тебя хранится token_hash, а не сам токен.
	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshTokenHash)
	if err != nil {
		return utils.NewUnauthorized("invalid refresh token")
	}

	currentAgent := r.UserAgent()
	if rt.UserAgent != currentAgent {
		return utils.NewUnauthorized("invalid refresh token")
	}

	// ШАГ 6️⃣ — проверяем срок действия
	if time.Now().After(rt.ExpiresAt) {
		_ = h.RefreshTokenRepo.Delete(r.Context(), refreshTokenHash)
		return utils.NewUnauthorized("refresh token expired")
	}

	// ШАГ 7️⃣ — ROTATION (очень важно 🔥), 👉 Удаляем старый refresh token
	// 📌 Почему:
	// 	если токен украдут — второй раз его не используют
	// 	это enterprise-security, uje ne nujen udalyaem po userId i po userAgent
	// if err := h.RefreshTokenRepo.Delete(r.Context(), refreshTokenHash); err != nil {
	// 	return utils.NewInternal("failed to rotate refresh token")
	// }

	// ШАГ 8️⃣ — создаём НОВЫЕ токены
	// Access token
	accessToken, err := utils.GenerateToken(rt.UserID)
	if err != nil {
		return utils.NewInternal("failed to generate access token")
	}
	// Refresh token
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return utils.NewInternal("failed to generate refresh token")
	}

	newRefreshHash := utils.HashToken(newRefreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// ШАГ 9️⃣ — сохраняем новый refresh token
	userAgent := r.UserAgent()
	// IP (корректно, по-профи), X-Forwarded-For нужен, если потом будешь за nginx / proxy
	ip := getClientIP(r)

	err = h.RefreshTokenRepo.DeleteByUserAndAgent(
		r.Context(),
		rt.UserID,
		userAgent,
	)
	if err != nil {
		return utils.NewInternal("failed to cleanup old refresh token")
	}

	err = h.RefreshTokenRepo.Create(
		r.Context(),
		rt.UserID,
		newRefreshHash,
		userAgent,
		ip,
		expiresAt,
	)
	if err != nil {
		return utils.NewInternal("failed to save refresh token")
	}

	// ШАГ 1. LOGIN — кладём refresh token в cookie
	// ❌ БЫЛО
	// utils.WriteJSON(w, http.StatusOK, map[string]any{
	// 	"access_token":  accessToken,
	// 	"refresh_token": newRefreshToken,
	// 	"expires_at":    expiresAt,
	// })

	// ✅ СТАНЕТ
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true esli HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
	// 📌 ПОЯСНЕНИЕ
	// 	HttpOnly → JS не видит
	// 	Path → cookie уходит ТОЛЬКО на refresh
	// 	SameSiteStrict → защита от CSRF
	// 	Secure → включишь когда будет HTTPS

	// 🔥 Ответ login теперь ТОЛЬКО access token
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
		"expires_at":   expiresAt,
	})

	return nil
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	// Читаем refresh_token из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		// Если куки нет, можно просто вернуть OK — пользователь и так "вышел"
		fmt.Println("No refresh cookie found")
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "logged out",
		})
		return nil
	}

	// Получаем userAgent
	userAgent := r.UserAgent()

	// Находим refresh token в БД и удаляем по userID + userAgent
	refreshHash := utils.HashToken(cookie.Value)
	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshHash)
	if err == nil {
		// Удаляем токен для текущего устройства
		_ = h.RefreshTokenRepo.DeleteByUserAndAgent(r.Context(), rt.UserID, userAgent)
	} else {
		fmt.Println("Refresh token not found in DB, maybe already deleted")
	}

	// Удаляем cookie у клиента
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	// Ответ
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
	return nil
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
	time.Sleep(2 * time.Second)
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

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"user_id":    rt.UserID,
		"user_agent": rt.UserAgent,
	})
	return nil
}
