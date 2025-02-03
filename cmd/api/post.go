package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/swarajroy/gophersocial/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// handler are not able to return errors??? why ?
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var postPayload CreatePostPayload

	if err := readJson(w, r, &postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := &store.Post{
		Title:   postPayload.Title,
		Content: postPayload.Content,
		Tags:    postPayload.Tags,
		// TODO after auth add correct userId
		UserID: 1,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetById(r.Context(), postId)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
