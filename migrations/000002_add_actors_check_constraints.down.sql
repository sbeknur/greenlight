ALTER TABLE actors DROP CONSTRAINT IF EXISTS actors_runtime_check;

ALTER TABLE actors DROP CONSTRAINT IF EXISTS actors_year_check;

ALTER TABLE actors DROP CONSTRAINT IF EXISTS genres_length_check;

ALTER TABLE actors DROP CONSTRAINT IF EXISTS actors_version_less10_check;