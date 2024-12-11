package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/edwrdc/digitally/internal/store"
	"github.com/go-chi/chi/v5"
)

type productKey string

const productCtx productKey = "product"

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

	if err := app.jsonResponse(w, http.StatusCreated, product); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {

	product := getProductFromContext(r)

	reviews, err := app.store.Reviews.GetByProductID(r.Context(), product.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	product.Reviews = reviews

	if err := app.jsonResponse(w, http.StatusOK, product); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "productID")
	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Products.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateProductPayload struct {
	Name        *string   `json:"name" validate:"omitempty,max=100"`
	Price       *float64  `json:"price" validate:"omitempty,number,gt=0"`
	Description *string   `json:"description" validate:"omitempty,max=1000"`
	Categories  *[]string `json:"categories" validate:"omitempty,min=1,max=5"`
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	product := getProductFromContext(r)

	var payload UpdateProductPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Name != nil {
		product.Name = *payload.Name
	}

	if payload.Price != nil {
		product.Price = *payload.Price
	}

	if payload.Description != nil {
		product.Description = *payload.Description
	}

	if payload.Categories != nil {
		product.Categories = *payload.Categories
	}

	if err := app.store.Products.Update(r.Context(), product); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		case errors.Is(err, store.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, product); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) productContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		productID := chi.URLParam(r, "productID")
		id, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		product, err := app.store.Products.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, productCtx, product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProductFromContext(r *http.Request) *store.Product {
	return r.Context().Value(productCtx).(*store.Product)
}
