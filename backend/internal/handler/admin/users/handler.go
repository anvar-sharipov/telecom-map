package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserRepo UserRepository
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return h.List(w, r)
	case http.MethodPost:
		return h.Create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	users, err := h.UserRepo.ListWithGroups(ctx)
	if err != nil {
		return err
	}

	result := make([]AdminUserDTO, 0, len(users))
	for _, u := range users {
		result = append(result, *ToAdminDTO(u))
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(result)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return nil
	}

	if len(req.Password) < 6 {
		http.Error(w, "password too short", http.StatusBadRequest)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: req.Username,
		FullName: req.FullName,
		Password: string(hash),
		IsActive: req.IsActive,
	}

	if err := h.UserRepo.CreateWithGroups(ctx, user, req.GroupIDs); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			http.Error(w, "username already exists", http.StatusConflict)
			return nil
		}
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(ToAdminDTO(user))
}
