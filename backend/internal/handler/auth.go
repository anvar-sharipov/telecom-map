package handler

import (
	"encoding/json"
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
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to hash password"})
		return
	}

	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	newUser := &domain.User{
		FullName: req.Fullname,
		Username: req.Username,
		Password: string(hashed),
		Phone:    phone,
		IsActive: true,
	}

	if err := h.UserRepo.Create(newUser); err != nil {

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
		"message": "User registered successfully",
		"token":   token,
	})
}
