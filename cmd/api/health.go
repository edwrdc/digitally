package main

import (
	"net/http"
)

// Healthcheck godoc
//
//	@Summary		API health check
//	@Description	Returns the current status of the API, environment, and version
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/healthz [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
