ALTER TABLE actors ADD CONSTRAINT actors_runtime_check CHECK (runtime >= 0);

ALTER TABLE actors ADD CONSTRAINT actors_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));

ALTER TABLE actors ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

ALTER TABLE actors ADD CONSTRAINT actors_version_less10_check CHECK (version < 10);