package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/swarajroy/gophersocial/internal/store"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("unauthorized"))
				return
			}

			// parse it -> base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("header 'Authorization' is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			// get credentials via DB or cache
			creds := strings.SplitN(string(decoded), ":", 2)
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("unauthorized"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("header 'Authorization' is malformed"))
			return
		}

		userID, err := app.auth.ValidateToken(parts[1])
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.getUser(ctx, userID)
		if err != nil {
			fmt.Println("app getUser failed")
			app.unauthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {

	if !app.config.cache.redisCfg.enabled {
		return app.store.Users.GetById(ctx, userID)
	}
	// Get from cache first -> read through cache
	user, err := app.cache.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("cache miss hitting db", "userID", userID)
		user, err = app.store.Users.GetById(ctx, userID)
		if err != nil {
			return nil, err
		}
		app.logger.Infow("inserting into the cache", "userID", userID, "user", user)
		err = app.cache.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
