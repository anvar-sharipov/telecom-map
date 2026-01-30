CREATE TABLE edges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_node UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    to_node UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    cable_type TEXT,
    length_m INTEGER,
    created_at TIMESTAMP DEFAULT now(),
    CHECK (from_node <> to_node)
);