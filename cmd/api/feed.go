package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feed, err := app.store.Products.GetUserFeed(ctx, int64(101))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
