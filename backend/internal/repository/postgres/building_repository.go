package postgres

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BuildingRepository struct {
	db *pgxpool.Pool
}

func NewBuildingRepository(db *pgxpool.Pool) repository.BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) Create(ctx context.Context, building *domain.Building) error {
	geoJSON, err := building.Geometry.ToGeoJSON()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO buildings (name, description, floors, geom)
		VALUES ($1, $2, $3, ST_SetSRID(ST_GeomFromGeoJSON($4), 4326))
	`

	_, err = r.db.Exec(
		ctx,
		query,
		building.Name,
		building.Description,
		building.Floors,
		geoJSON,
	)

	return err
}
