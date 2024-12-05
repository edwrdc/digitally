package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/edwrdc/digitally/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateProductPayload struct {
	Name        string   `json:"name" validate:"required,max=100"`
	Price       float64  `json:"price" validate:"required,number,gt=0"`
	Description string   `json:"description" validate:"required,max=1000"`
	Categories  []string `json:"categories" validate:"required,min=1,max=5"`
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateProductPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &store.Product{
		UserID:      1,
		Name:        payload.Name,
		Price:       payload.Price,
		Description: payload.Description,
		Categories:  payload.Categories,
	}
	ctx := r.Context()

	if err := app.store.Products.Create(ctx, product); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, product); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "productID")
	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	product, err := app.store.Products.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
