package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/swarajroy/gophersocial/internal/auth"
	"github.com/swarajroy/gophersocial/internal/store"
	"github.com/swarajroy/gophersocial/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	cfg := config{}
	store := store.NewMockDbStorage()
	cache := cache.NewMockCacheStorage()
	auth := auth.NewTestAuthenticator()
	app := &application{
		config: cfg,
		logger: logger,
		store:  store,
		cache:  cache,
		auth:   auth,
	}

	return app
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}
