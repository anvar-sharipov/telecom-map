package domain

import (
	"time"

	"github.com/google/uuid"
)

type NodeType string

const (
	NodeTelecom  NodeType = "telecom"
	NodeSplice   NodeType = "splice"
	NodeBuilding NodeType = "building"
)

type Node struct {
	ID        uuid.UUID
	Name      string
	Type      NodeType
	Lon       float64
	Lat       float64
	CreatedAt time.Time
}
