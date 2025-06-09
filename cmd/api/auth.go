package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swarajroy/gophersocial/internal/mailer"
	"github.com/swarajroy/gophersocial/internal/store"
)

var (
	ErrRequiresHigherPrivelege = errors.New("requires higher privelege")
	ErrInvalidRole             = errors.New("valid roles are admin | user | moderator")
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=5,max=10"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=15"`
	Role     string `json:"role" validate:"required,min=4,max=20"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

func (app *application) postAuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	role := store.Role{
		Name: payload.Role,
	}

	_, err := app.store.Roles.GetRoleByName(r.Context(), role.Name)
	if err != nil {
		app.badRequestError(w, r, ErrInvalidRole)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Role:     role,
	}

	//hash
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	//create and invite
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	ctx := r.Context()
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.email.expiry); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	//send email
	isProdEnv := app.config.env == "production"
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	_, err = app.mailer.Send(ctx, mailer.UserInvitationTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback the user invitation and the user record
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting the user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type CreateTokenPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=15"`
}

func (app *application) postTokenHandler(w http.ResponseWriter, r *http.Request) {
	//parse payload creds
	var payload CreateTokenPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	//fetch the user from the creds
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.IsValid(payload.Password); err != nil {
		app.unauthorizedError(w, r, err)
		return
	}
	//generate the token -> add claims
	token, err := app.auth.GenerateToken(user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	//send it to the client
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if the user is the one doing something with the post
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// check for role precedence
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r, ErrRequiresHigherPrivelege)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, requiredRole string) (bool, error) {

	role, err := app.store.Roles.GetRoleByName(ctx, requiredRole)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
