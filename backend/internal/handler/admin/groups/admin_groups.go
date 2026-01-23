package groups

import (
	"encoding/json"
	"net/http"

	"github.com/anvar-sharipov/telecom-map/internal/repository"
)

type AdminGroupHandler struct {
	GroupRepo *repository.GroupRepository
}

type GroupDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *AdminGroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	isActiveParam := r.URL.Query().Get("is_active")

	var (
		groups []repository.Group
		err    error
	)

	if isActiveParam == "" {
		groups, err = h.GroupRepo.ListAll(ctx)
	} else {
		isActive := isActiveParam == "true"
		groups, err = h.GroupRepo.ListByActive(ctx, isActive)
	}

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(groups)
}
