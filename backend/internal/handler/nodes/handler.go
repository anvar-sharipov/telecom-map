package nodes

import (
	"encoding/json"
	"net/http"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
)

type Handler struct {
	NodeRepo repository.NodeRepository
}

type createNodeRequest struct {
	Name string          `json:"name"`
	Type domain.NodeType `json:"type"`
	Lon  float64         `json:"lon"`
	Lat  float64         `json:"lat"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	var req createNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	node := &domain.Node{
		Name: req.Name,
		Type: req.Type,
		Lon:  req.Lon,
		Lat:  req.Lat,
	}

	err := h.NodeRepo.Create(r.Context(), node)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(node)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) error {
	nodes, err := h.NodeRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(nodes)
}
