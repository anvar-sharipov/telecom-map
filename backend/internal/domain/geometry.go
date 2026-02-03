package domain

import (
	"encoding/json"
)

type Geometry struct {
	Type        string        // "Polygon"
	Coordinates [][][]float64 // координаты полигона
}

func (g Geometry) ToGeoJSON() (string, error) {
	data := map[string]interface{}{
		"type":        g.Type,
		"coordinates": g.Coordinates,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
