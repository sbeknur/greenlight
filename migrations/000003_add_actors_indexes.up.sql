CREATE INDEX IF NOT EXISTS actors_name_idx ON actors USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS actors_surname_idx ON actors USING GIN (to_tsvector('simple', surname));