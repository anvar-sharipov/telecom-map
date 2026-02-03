package handler

import (
	"encoding/json"
	"net/http"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/anvar-sharipov/telecom-map/internal/utils"
)

type BuildingHandler struct {
	BuildingRepo repository.BuildingRepository
}

func (h *BuildingHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.NewMethodNotAllowed("method not allowed")
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Floors      int    `json:"floors"`
		Geometry    struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.NewBadRequest("invalid request body")
	}

	// 🔒 базовая валидация
	if req.Geometry.Type != "Polygon" {
		return utils.NewBadRequest("geometry must be Polygon")
	}

	if len(req.Geometry.Coordinates) == 0 {
		return utils.NewBadRequest("coordinates required")
	}

	building := &domain.Building{
		Name:        req.Name,
		Description: req.Description,
		Floors:      req.Floors,
		Geometry: domain.Geometry{
			Type:        req.Geometry.Type,
			Coordinates: req.Geometry.Coordinates,
		},
	}

	if err := h.BuildingRepo.Create(r.Context(), building); err != nil {
		return utils.NewInternal("failed to create building")
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "building created",
	})

	return nil
}
