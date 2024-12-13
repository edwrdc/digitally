package main

import (
	"net/http"

	"github.com/edwrdc/digitally/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get user's product feed
//	@Description	Retrieves a paginated feed of products for the user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int		false	"Number of items per page"	default(20)
//	@Param			offset		query		int		false	"Offset for pagination"		default(0)
//	@Param			sort		query		string	false	"Sort order (asc/desc)"		default(desc)
//	@Param			category	query		string	false	"Category to filter by"
//	@Param			search		query		string	false	"Search term"
//	@Param			since		query		string	false	"Since date (YYYY-MM-DD)"
//	@Param			until		query		string	false	"Until date (YYYY-MM-DD)"
//	@Success		200			{array}		[]store.UserFeedProduct
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginationFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Products.GetUserFeed(ctx, int64(101), fq)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
