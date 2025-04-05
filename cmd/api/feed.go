package main

import (
	"net/http"

	"github.com/swarajroy/gophersocial/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg := store.PaginatedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	pg, err := pg.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(pg); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(8), pg)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
