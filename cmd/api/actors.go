package main

import (
	"fmt"
	"net/http"

	"github.com/sbeknur/greenlight/internal/data"
)

func (app *application) createActorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID        int64     `json:"id"`
		Name     string    `json:"title"`
		Age      int32     `json:"year,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	actor := &data.Actor{
		Name:   input.Name,
		Age:    input.Age,
	}

	err = app.models.Actors.Insert(actor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/actors/%d", actor.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"actor": actor}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}