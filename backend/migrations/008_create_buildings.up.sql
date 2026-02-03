CREATE TABLE buildings (
    id SERIAL PRIMARY KEY,

    -- геометрия здания
    geom geometry(Polygon, 4326) NOT NULL,

    -- метаданные
    name TEXT,
    floors INTEGER,

    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_buildings_geom
ON buildings
USING GIST (geom);