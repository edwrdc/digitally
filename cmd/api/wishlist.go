package main

import (
	"net/http"

	"github.com/edwrdc/digitally/internal/store"
)

// AddToWishlist godoc
//
//	@Summary		Add product to wishlist
//	@Description	Adds a product to the user's wishlist
//	@Tags			wishlist
//	@Accept			json
//	@Produce		json
//	@Param			productID	path		int	true	"Product ID"
//	@Success		204			{object}	nil
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error	"Product not found"
//	@Failure		409			{object}	error	"Product already in wishlist"
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist/{productID} [put]
func (app *application) addProductToWishlistHandler(w http.ResponseWriter, r *http.Request) {

	product := getProductFromContext(r)
	user := getUserFromContext(r)

	ctx := r.Context()

	if err := app.store.Wishlist.Add(ctx, user.ID, product.ID); err != nil {
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

// RemoveFromWishlist godoc
//
//	@Summary		Remove product from wishlist
//	@Description	Removes a product from the user's wishlist
//	@Tags			wishlist
//	@Accept			json
//	@Produce		json
//	@Param			productID	path		int	true	"Product ID"
//	@Success		204			{object}	nil
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error	"Product not found"
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist/{productID} [delete]
func (app *application) removeProductFromWishlistHandler(w http.ResponseWriter, r *http.Request) {

	product := getProductFromContext(r)

	user := getUserFromContext(r)

	ctx := r.Context()

	if err := app.store.Wishlist.Remove(ctx, user.ID, product.ID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
