package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		// Healthcheck
		r.Get("/healthz", app.healthcheckHandler)

		docsURL := fmt.Sprintf(":%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// Products
		r.Route("/products", func(r chi.Router) {

			r.Post("/", app.createProductHandler)

			r.Route("/{productID}", func(r chi.Router) {
				r.Use(app.productContextMiddleware)
				r.Get("/", app.getProductHandler)

				r.Delete("/", app.deleteProductHandler)
				r.Patch("/", app.updateProductHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/wishlist", func(r chi.Router) {
			r.Route("/{productID}", func(r chi.Router) {
				r.Use(app.productContextMiddleware)
				r.Put("/", app.addProductToWishlistHandler)
				r.Delete("/", app.removeProductFromWishlistHandler)
			})
		})
	})

	return r
}
