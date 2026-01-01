package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
	"github.com/anvar-sharipov/telecom-map/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo *postgres.UserRepository
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		Fullname        string `json:"fullname"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Password != req.ConfirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "passwords do not match",
		})
		return
	}

	if req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "password cant be empty",
		})
		return
	}

	if req.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "username cant be empty",
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to hash password"})
		return
	}

	newUser := &domain.User{
		FullName: req.Fullname,
		Username: req.Username,
		Password: string(hashed),
		IsActive: true,
	}

	if err := h.UserRepo.Create(newUser); err != nil {
		fmt.Println("UserRepo.Create error:", err) // <--- вот здесь
		w.Header().Set("Content-Type", "application/json")

		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "username already exists",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}

	token, err := utils.GenerateToken(newUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
		"token":   token,
	})
}
