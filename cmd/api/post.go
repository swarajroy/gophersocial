package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/swarajroy/gophersocial/internal/store"
)

type postKey string

const (
	postCtx postKey = "post"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required=true,max=100"`
	Content string   `json:"content" validate:"required=true,max=200"`
	Tags    []string `json:"tags"`
}

// handler are not able to return errors??? why ?
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var postPayload CreatePostPayload

	if err := readJson(w, r, &postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(postPayload); err != nil {
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

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetPostById(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	if err := app.store.Posts.Delete(r.Context(), post.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	log.Printf("payload = %+v", payload)

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		ctx := r.Context()
		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
