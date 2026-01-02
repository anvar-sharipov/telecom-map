package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
	"github.com/anvar-sharipov/telecom-map/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo         *postgres.UserRepository
	RefreshTokenRepo *repository.RefreshTokenRepository
}

// ---------------- REGISTER ----------------
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
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

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return utils.NewInternal("failed to generate token")
	}

	// refreshToken, err := utils.GenerateRefreshToken()
	// if err != nil {
	// 	return utils.NewInternal("failed to generate refresh token")
	// }

	// expiresAt := time.Now().Add(7 * 24 * time.Hour)
	// err = h.RefreshTokenRepo.Create(
	// 	context.Background(),
	// 	user.ID,
	// 	refreshToken,
	// 	expiresAt,
	// )
	// if err != nil {
	// 	return utils.NewInternal("failed to save refresh token")
	// }

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "login successful",
		"token":   token,
	})
	return nil
}
