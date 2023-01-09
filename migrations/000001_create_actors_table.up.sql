CREATE TABLE IF NOT EXISTS actors (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
	surname text NOT NULL,
    age integer NOT NULL
);