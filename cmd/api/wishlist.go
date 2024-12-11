package main

import (
	"net/http"

	"github.com/edwrdc/digitally/internal/store"
)

type WishlistRequest struct {
	UserID int64 `json:"user_id"`
}

func (app *application) getWishlistHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) addProductToWishlistHandler(w http.ResponseWriter, r *http.Request) {

	product := getProductFromContext(r)

	// TODO: will do auth later
	var payload WishlistRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Wishlist.Add(ctx, payload.UserID, product.ID); err != nil {
		switch {
		case err == store.ErrConflict:
			app.conflictResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) removeProductFromWishlistHandler(w http.ResponseWriter, r *http.Request) {

	product := getProductFromContext(r)

	// TODO: will do auth later
	var payload WishlistRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Wishlist.Remove(ctx, payload.UserID, product.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
