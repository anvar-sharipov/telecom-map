package admin

import (
	"encoding/json"
	"net/http"

	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
)

type AdminHandler struct {
	UserRepo *postgres.UserRepository
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	users, err := h.UserRepo.ListWithGroups(ctx)
	if err != nil {
		return err
	}

	result := make([]UserDTO, 0, len(users))
	for _, u := range users {
		result = append(result, mapUserToDTO(u))
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(result)
}
