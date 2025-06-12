package main

import (
	"context"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {

	app := newTestApplication(t)
	mux := app.mount()
	testToken := "abc123"
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected %d got %d", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("should get user", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)

		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, rr.Code)
		}

	})
}
