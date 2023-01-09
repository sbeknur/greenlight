package data

import (
	"context"
	"errors"

	"fmt"
	"time"

	"database/sql"
)

type Actor struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
	Age     int32  `json:"age,omitempty"`
}

type ActorModel struct {
	DB *sql.DB
}

func (a ActorModel) Insert(actor *Actor) error {
	query :=
		`INSERT INTO actors (name, surname, age)
		 VALUES ($1, $2, $3)
		 RETURNING id`
	args := []any{actor.Name, actor.Surname, actor.Age}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ID)
}

func (a ActorModel) Get(id int64) (*Actor, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, surname, age
		FROM actors
		WHERE id = $1`
	var actor Actor
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, id).Scan(
		&actor.ID,
		&actor.Name,
		&actor.Surname,
		&actor.Age,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &actor, nil
}

func (a ActorModel) GetAll(name string, surname string, filters Filters) ([]*Actor, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, surname, age
		FROM actors
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', surname) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{name, surname, filters.limit(), filters.offset()}
	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	actors := []*Actor{}
	for rows.Next() {
		var actor Actor
		err := rows.Scan(
			&totalRecords,
			&actor.ID,
			&actor.Name,
			&actor.Surname,
			&actor.Age,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		actors = append(actors, &actor)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return actors, metadata, nil
}
