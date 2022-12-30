package data

import (
	"database/sql"
)

type Actor struct {
	ID        int64     `json:"id"`
	Name     string    `json:"title"`
	Age      int32     `json:"year,omitempty"`
}

type ActorModel struct {
	DB *sql.DB
}

func (m ActorModel) Insert(actor *Actor) error {
	query :=
		`INSERT INTO actors (name, age)
		 VALUES ($1, $2)
		 RETURNING id`
	args := []any{actor.Name, actor.Age}

	return m.DB.QueryRow(query, args...).Scan(&actor.ID)
}