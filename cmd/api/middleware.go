package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/edwrdc/digitally/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicAuthResponse(w, r, fmt.Errorf("no auth header"))
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
				app.unauthorizedBasicAuthResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicAuthResponse(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass
			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 || credentials[0] != username || credentials[1] != pass {
				app.unauthorizedBasicAuthResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == " " {
			app.unauthorizedResponse(w, r, fmt.Errorf("no auth header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			app.unauthorizedResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.getUserWithCache(ctx, userID)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkProductOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := getUserFromContext(r)
		product := getProductFromContext(r)

		if product.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) getUserWithCache(ctx context.Context, userID int64) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByID(ctx, userID)
	}

	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if err = app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
