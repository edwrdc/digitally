package main

import (
	"log"
	"net/http"
)

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: Introduce structured logging later
	log.Printf("Internal server error: %s, path:%s, error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "encountered an error while processing the request")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %s, path:%s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found: %s, path:%s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("Edit conflict: %s, path:%s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusConflict, "edit conflict")
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) {
	type envelope struct {
		Data any `json:"data"`
	}

	writeJSON(w, status, &envelope{Data: data})
}
