package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/edwrdc/digitally/internal/store"
	"github.com/go-chi/chi/v5"
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

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "productID")
	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	product, err := app.store.Products.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, product); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
