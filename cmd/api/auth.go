package main

import (
	"net/http"

	"github.com/swarajroy/gophersocial/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=5,max=10"`
	Email    string `json:"email" validate:"required,email,min=3,max=10"`
	Password string `json:"password" validate:"required,min-3,max=15"`
}

func (app *application) postAuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJson(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	//hash
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	//create and invite
	ctx := r.Context()
	if err := app.store.Users.CreateAndInvite(ctx, user, "uuid4"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	//send email
	if err := writeJson(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
