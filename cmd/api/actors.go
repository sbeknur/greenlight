package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sbeknur/greenlight/internal/data"
	"github.com/sbeknur/greenlight/internal/validator"
)

func (app *application) createActorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Age     int32  `json:"age"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	actor := &data.Actor{
		Name:    input.Name,
		Surname: input.Surname,
		Age:     input.Age,
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

func (app *application) showActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	actors, err := app.models.Actors.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"actor": actors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


func (app *application) listActorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Surname  string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")
	
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "name", "surname", "age", "-id", "-name", "surname", "-age"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	actors, metadata, err := app.models.Actors.GetAll(input.Name, input.Surname, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"actors": actors, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
