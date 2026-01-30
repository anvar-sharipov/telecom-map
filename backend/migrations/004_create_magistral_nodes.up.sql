CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT,
    type TEXT CHECK (type IN ('telecom', 'splice', 'building')),
    lon DOUBLE PRECISION NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);