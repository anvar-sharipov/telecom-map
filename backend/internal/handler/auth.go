// package handler

// import (
// 	"encoding/json"
// 	"fmt"

// 	// "net"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/anvar-sharipov/telecom-map/internal/domain"
// 	"github.com/anvar-sharipov/telecom-map/internal/repository"
// 	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
// 	"github.com/anvar-sharipov/telecom-map/internal/service"
// 	"github.com/anvar-sharipov/telecom-map/internal/utils"
// 	"golang.org/x/crypto/bcrypt"
// )

// // func getClientIP(r *http.Request) string {
// // 	// 1. Если есть proxy
// // 	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
// // 		// XFF может быть "client, proxy1, proxy2"
// // 		return strings.TrimSpace(strings.Split(xff, ",")[0])
// // 	}

// // 	// 2. Без proxy
// // 	host, _, err := net.SplitHostPort(r.RemoteAddr)
// // 	if err != nil {
// // 		return r.RemoteAddr
// // 	}

// // 	return host
// // }

// type AuthHandler struct {
// 	UserRepo         *postgres.UserRepository
// 	RefreshTokenRepo *repository.RefreshTokenRepository
// 	AuthService      *service.AuthService
// }

// // ---------------- REGISTER ----------------
// func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method != http.MethodPost {
// 		return utils.NewMethodNotAllowed("method not allowed")
// 	}

// 	var req struct {
// 		Fullname        string `json:"fullname"`
// 		Username        string `json:"username"`
// 		Password        string `json:"password"`
// 		ConfirmPassword string `json:"confirm_password"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return utils.NewBadRequest("Invalid request body")
// 	}

// 	if req.Password != req.ConfirmPassword {
// 		return utils.NewBadRequest("passwords do not match")
// 	}

// 	if req.Password == "" {
// 		return utils.NewBadRequest("password cant be empty")
// 	}

// 	if req.Username == "" {
// 		return utils.NewBadRequest("username cant be empty")
// 	}

// 	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return utils.NewInternal("Failed to hash password")
// 	}

// 	newUser := &domain.User{
// 		FullName: req.Fullname,
// 		Username: req.Username,
// 		Password: string(hashed),
// 		IsActive: true,
// 	}

// 	if err := h.UserRepo.Create(newUser); err != nil {
// 		if strings.Contains(err.Error(), "duplicate key") {
// 			return utils.NewConflict("username already exists")
// 		}
// 		return utils.NewInternal("internal server error")
// 	}

// 	// 🔥 АВТОЛОГИН — используем сервис (claudeai)
// 	tokens, err := h.AuthService.CreateSession(r.Context(), service.SessionData{
// 		UserID:    newUser.ID,
// 		UserAgent: r.UserAgent(),
// 		IP:        utils.GetClientIP(r),
// 	})
// 	if err != nil {
// 		return utils.NewInternal("failed to create session")
// 	}

// 	// Устанавливаем cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    tokens.RefreshToken,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   false, // true при HTTPS
// 		SameSite: http.SameSiteLaxMode,
// 		Expires:  tokens.ExpiresAt,
// 	})

// 	utils.WriteJSON(w, http.StatusCreated, map[string]any{
// 		"message":      "user registered successfully",
// 		"access_token": tokens.AccessToken,
// 		"expires_at":   tokens.ExpiresAt,
// 	})

// 	return nil

// 	// // ✅ АВТОЛОГИН после регистрации - ДЕЛАЕМ ТО ЖЕ, ЧТО И В LOGIN
// 	// // prawilnyy wariant no dubliruetsya kod
// 	// // 1️⃣ ACCESS TOKEN
// 	// accessToken, err := utils.GenerateToken(newUser.ID)
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to generate access token")
// 	// }

// 	// // 2️⃣ REFRESH TOKEN
// 	// refreshToken, err := utils.GenerateRefreshToken()
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to generate refresh token")
// 	// }

// 	// // 3️⃣ HASH refresh token
// 	// refreshTokenHash := utils.HashToken(refreshToken)

// 	// // 4️⃣ EXPIRE: 7 дней или 30 секунд для теста
// 	// expiresAt := time.Now().Add(30 * time.Second) // Для теста

// 	// // 5️⃣ SAVE refresh token
// 	// userAgent := r.UserAgent()
// 	// ip := utils.GetClientIP(r)

// 	// // Удаляем старый токен (на всякий случай)
// 	// err = h.RefreshTokenRepo.DeleteByUserAndAgent(
// 	// 	r.Context(),
// 	// 	newUser.ID,
// 	// 	userAgent,
// 	// )
// 	// if err != nil {
// 	// 	// Можно логировать, но не прерывать
// 	// 	log.Println("failed to cleanup old refresh token on register:", err)
// 	// }

// 	// // Сохраняем новый
// 	// err = h.RefreshTokenRepo.Create(
// 	// 	r.Context(),
// 	// 	newUser.ID,
// 	// 	refreshTokenHash,
// 	// 	userAgent,
// 	// 	ip,
// 	// 	expiresAt,
// 	// )
// 	// if err != nil {
// 	// 	log.Println("refresh token save error on register:", err)
// 	// 	return utils.NewInternal("failed to save refresh token")
// 	// }

// 	// // 6️⃣ УСТАНАВЛИВАЕМ COOKIE С REFRESH TOKEN
// 	// http.SetCookie(w, &http.Cookie{
// 	// 	Name:     "refresh_token",
// 	// 	Value:    refreshToken,
// 	// 	Path:     "/",
// 	// 	HttpOnly: true,
// 	// 	Secure:   false, // true при HTTPS
// 	// 	SameSite: http.SameSiteLaxMode,
// 	// 	Expires:  expiresAt,
// 	// })

// 	// // 7️⃣ ВОЗВРАЩАЕМ ТОЛЬКО ACCESS TOKEN (refresh уже в cookie)
// 	// utils.WriteJSON(w, http.StatusCreated, map[string]any{
// 	// 	"message":      "user registered successfully",
// 	// 	"access_token": accessToken,
// 	// 	"expires_at":   expiresAt,
// 	// })
// 	// return nil
// }

// // ---------------- LOGIN ----------------
// func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method != http.MethodPost {
// 		return utils.NewMethodNotAllowed("method not allowed")
// 	}

// 	var req struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return utils.NewBadRequest("Invalid request body")
// 	}

// 	if req.Username == "" {
// 		return utils.NewBadRequest("username cant be empty")
// 	}

// 	if req.Password == "" {
// 		return utils.NewBadRequest("password cant be empty")
// 	}

// 	user, err := h.UserRepo.GetByUsername(req.Username)
// 	if err != nil {
// 		return utils.NewUnauthorized("invalid credentials")
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
// 		return utils.NewUnauthorized("invalid credentials")
// 	}

// 	// 🔥 СОЗДАЁМ СЕССИЮ через сервис (claideai)
// 	tokens, err := h.AuthService.CreateSession(r.Context(), service.SessionData{
// 		UserID:    user.ID,
// 		UserAgent: r.UserAgent(),
// 		IP:        utils.GetClientIP(r),
// 	})
// 	if err != nil {
// 		return utils.NewInternal("failed to create session")
// 	}

// 	// Устанавливаем cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    tokens.RefreshToken,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   false, // true при HTTPS
// 		SameSite: http.SameSiteLaxMode,
// 		Expires:  tokens.ExpiresAt,
// 	})

// 	utils.WriteJSON(w, http.StatusOK, map[string]any{
// 		"access_token": tokens.AccessToken,
// 		"expires_at":   tokens.ExpiresAt,
// 	})

// 	return nil

// 	// // 1️⃣ ACCESS TOKEN rabotaet no bez service
// 	// accessToken, err := utils.GenerateToken(user.ID)
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to generate access token")
// 	// }

// 	// // 2️⃣ REFRESH TOKEN
// 	// refreshToken, err := utils.GenerateRefreshToken()
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to generate refresh token")
// 	// }

// 	// // 3️⃣ HASH refresh token
// 	// refreshTokenHash := utils.HashToken(refreshToken)

// 	// // 4️⃣ EXPIRE: срок действия, текущий момент времени + 7 дней
// 	// // expiresAt := time.Now().Add(7 * 24 * time.Hour)
// 	// expiresAt := time.Now().Add(30 * time.Second) // 🔁 Refresh token — 30 секунд

// 	// // 5️⃣ SAVE refresh token
// 	// userAgent := r.UserAgent()

// 	// ip := utils.GetClientIP(r)

// 	// err = h.RefreshTokenRepo.DeleteByUserAndAgent(
// 	// 	r.Context(),
// 	// 	user.ID,
// 	// 	userAgent,
// 	// )
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to cleanup old refresh token")
// 	// }

// 	// err = h.RefreshTokenRepo.Create(
// 	// 	r.Context(),
// 	// 	user.ID,
// 	// 	refreshTokenHash,
// 	// 	userAgent,
// 	// 	ip,
// 	// 	expiresAt,
// 	// )
// 	// if err != nil {
// 	// 	log.Println("refresh token save error:", err)
// 	// 	return utils.NewInternal("failed to save refresh token")
// 	// }

// 	// // ❌ БЫЛО
// 	// // 6️⃣ RESPONSE
// 	// // utils.WriteJSON(w, http.StatusOK, map[string]any{
// 	// // 	"access_token":  accessToken,
// 	// // 	// "refresh_token": refreshToken,
// 	// // 	"expires_at":    expiresAt,
// 	// // })

// 	// // ✅ СТАНЕТ
// 	// http.SetCookie(w, &http.Cookie{
// 	// 	Name:     "refresh_token",
// 	// 	Value:    refreshToken,
// 	// 	Path:     "/",
// 	// 	HttpOnly: true,
// 	// 	Secure:   false, // true при HTTPS
// 	// 	SameSite: http.SameSiteLaxMode,
// 	// 	Expires:  expiresAt,
// 	// })
// 	// // SameSite — Strict ОК, но знай нюанс
// 	// // ✔ Для SPA + same domain — отлично
// 	// // 	⚠ Если потом будет:
// 	// // 	фронт на другом домене
// 	// // 	мобильное приложение
// 	// // 	OAuth
// 	// // 	👉 тогда меняют на:
// 	// // 	SameSite: http.SameSiteLaxMode
// 	// // 	Пока оставляем Strict, ты сделал правильно.

// 	// utils.WriteJSON(w, http.StatusOK, map[string]any{
// 	// 	"access_token": accessToken,
// 	// 	"expires_at":   expiresAt,
// 	// })
// 	// return nil
// }

// func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method != http.MethodPost {
// 		return utils.NewMethodNotAllowed("method not allowed")
// 	}

// 	// ✅ СТАНЕТ
// 	cookie, err := r.Cookie("refresh_token")
// 	if err != nil {
// 		return utils.NewUnauthorized("refresh token missing")
// 	}

// 	// 📌 Body больше не нужен вообще
// 	refreshToken := cookie.Value

// 	// ШАГ 4️⃣ — хешируем refresh token
// 	refreshTokenHash := utils.HashToken(refreshToken)

// 	// ШАГ 5️⃣ — ищем refresh token в БД, ❗ В БД у тебя хранится token_hash, а не сам токен.
// 	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshTokenHash)
// 	if err != nil {
// 		return utils.NewUnauthorized("invalid refresh token")
// 	}

// 	currentAgent := r.UserAgent()
// 	if rt.UserAgent != currentAgent {
// 		return utils.NewUnauthorized("invalid refresh token")
// 	}

// 	// ШАГ 6️⃣ — проверяем срок действия
// 	if time.Now().After(rt.ExpiresAt) {
// 		_ = h.RefreshTokenRepo.Delete(r.Context(), refreshTokenHash)
// 		return utils.NewUnauthorized("refresh token expired")
// 	}

// 	// 🔥 РОТАЦИЯ СЕССИИ через сервис (claudeai)
// 	tokens, err := h.AuthService.RotateSession(r.Context(), service.SessionData{
// 		UserID:    rt.UserID,
// 		UserAgent: r.UserAgent(),
// 		IP:        utils.GetClientIP(r),
// 	})
// 	if err != nil {
// 		return utils.NewInternal("failed to rotate session")
// 	}

// 	// Устанавливаем новый cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    tokens.RefreshToken,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   false,
// 		SameSite: http.SameSiteLaxMode,
// 		Expires:  tokens.ExpiresAt,
// 	})

// 	utils.WriteJSON(w, http.StatusOK, map[string]any{
// 		"access_token": tokens.AccessToken,
// 		"expires_at":   tokens.ExpiresAt,
// 	})

// 	return nil

// 	// // ШАГ 8️⃣ — создаём НОВЫЕ токены rabotaet no staroe bez service
// 	// // Access token
// 	// accessToken, err := utils.GenerateToken(rt.UserID)
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to generate access token")
// 	// }
// 	// // Refresh token
// 	// newRefreshToken, err := utils.GenerateRefreshToken()
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to generate refresh token")
// 	// }

// 	// newRefreshHash := utils.HashToken(newRefreshToken)
// 	// // expiresAt := time.Now().Add(7 * 24 * time.Hour)
// 	// expiresAt := time.Now().Add(30 * time.Second) // 🔁 Refresh token — 30 секунд

// 	// // ШАГ 9️⃣ — сохраняем новый refresh token
// 	// userAgent := r.UserAgent()
// 	// ip := utils.GetClientIP(r)

// 	// err = h.RefreshTokenRepo.DeleteByUserAndAgent(
// 	// 	r.Context(),
// 	// 	rt.UserID,
// 	// 	userAgent,
// 	// )
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to cleanup old refresh token")
// 	// }

// 	// err = h.RefreshTokenRepo.Create(
// 	// 	r.Context(),
// 	// 	rt.UserID,
// 	// 	newRefreshHash,
// 	// 	userAgent,
// 	// 	ip,
// 	// 	expiresAt,
// 	// )
// 	// if err != nil {
// 	// 	return utils.NewInternal("failed to save refresh token")
// 	// }

// 	// // ✅ СТАНЕТ
// 	// http.SetCookie(w, &http.Cookie{
// 	// 	Name:     "refresh_token",
// 	// 	Value:    newRefreshToken,
// 	// 	Path:     "/",
// 	// 	HttpOnly: true,
// 	// 	Secure:   false, // true esli HTTPS
// 	// 	SameSite: http.SameSiteLaxMode,
// 	// 	Expires:  expiresAt,
// 	// })

// 	// // 🔥 Ответ login теперь ТОЛЬКО access token
// 	// utils.WriteJSON(w, http.StatusOK, map[string]any{
// 	// 	"access_token": accessToken,
// 	// 	"expires_at":   expiresAt,
// 	// })

// 	// return nil
// }

// func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method != http.MethodPost {
// 		return utils.NewMethodNotAllowed("method not allowed")
// 	}

// 	// Читаем refresh_token из cookie
// 	cookie, err := r.Cookie("refresh_token")
// 	if err != nil {
// 		// Если куки нет, можно просто вернуть OK — пользователь и так "вышел"
// 		fmt.Println("No refresh cookie found")
// 		utils.WriteJSON(w, http.StatusOK, map[string]string{
// 			"message": "logged out",
// 		})
// 		return nil
// 	}

// 	// Получаем userAgent
// 	// userAgent := r.UserAgent()

// 	// Находим refresh token в БД и удаляем по userID + userAgent
// 	refreshHash := utils.HashToken(cookie.Value)
// 	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshHash)
// 	if err == nil {
// 		// Удаляем токен для текущего устройства
// 		// _ = h.RefreshTokenRepo.DeleteByUserAndAgent(r.Context(), rt.UserID, userAgent)
// 		// 🔥 ИНВАЛИДАЦИЯ СЕССИИ через сервис (claudeAi)
// 		_ = h.AuthService.InvalidateSession(r.Context(), rt.UserID, r.UserAgent())
// 	} else {
// 		fmt.Println("Refresh token not found in DB, maybe already deleted")
// 	}

// 	// Удаляем cookie у клиента
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    "",
// 		Path:     "/",
// 		HttpOnly: true,
// 		Expires:  time.Unix(0, 0),
// 		MaxAge:   -1,
// 	})

// 	// Ответ
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{
// 		"message": "logged out",
// 	})
// 	return nil
// }

// func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
// 	if r.Method != http.MethodGet {
// 		return utils.NewMethodNotAllowed("method not allowed")
// 	}

// 	cookie, err := r.Cookie("refresh_token")
// 	if err != nil {
// 		return utils.NewUnauthorized("not logged in")
// 	}

// 	refreshHash := utils.HashToken(cookie.Value)
// 	rt, err := h.RefreshTokenRepo.GetByTokenHash(r.Context(), refreshHash)
// 	if err != nil {
// 		return utils.NewUnauthorized("not logged in")
// 	}

// 	utils.WriteJSON(w, http.StatusOK, map[string]any{
// 		"user_id":    rt.UserID,
// 		"user_agent": rt.UserAgent,
// 	})
// 	return nil
// }

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
