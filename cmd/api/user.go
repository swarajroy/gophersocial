package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/swarajroy/gophersocial/internal/store"
)

type userKey string

const (
	userCtx userKey = "user"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil || userId < 1 {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.getUser(r.Context(), userId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type FollowerPayload struct {
	UserId int64 `json:"user_Id"`
}

func (app *application) putFollowHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	var follower FollowerPayload

	if err := readJson(w, r, &follower); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := app.store.Followers.Follow(r.Context(), user.ID, follower.UserId); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) putUnfollowHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var unFollower FollowerPayload

	if err := readJson(w, r, &unFollower); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), user.ID, unFollower.UserId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

// TOOO swagger docs
func (app *application) putActivateUser(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
