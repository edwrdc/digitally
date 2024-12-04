package main

import (
	"net/http"

	"github.com/edwrdc/digitally/internal/store"
)

type CreateProductPayload struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateProductPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	product := &store.Product{
		UserID:      1,
		Name:        payload.Name,
		Price:       payload.Price,
		Description: payload.Description,
		Categories:  payload.Categories,
	}

	if err := app.store.Products.Create(ctx, product); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, product); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
