package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2" // http-swagger middleware
	"github.com/swarajroy/gophersocial/docs"
	"github.com/swarajroy/gophersocial/internal/auth"
	"github.com/swarajroy/gophersocial/internal/mailer"
	"github.com/swarajroy/gophersocial/internal/store"
	"go.uber.org/zap"
)

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
	auth   auth.TokenGenerator
}

type authConfig struct {
	basic basicAuthConfig
	jwt   jwtConfig
}

type jwtConfig struct {
	secret string
	host   string
	exp    time.Duration
}

type basicAuthConfig struct {
	user string
	pass string
}

type config struct {
	addr        string
	dbConfig    dbConfig
	env         string
	apiURL      string
	email       emailConfig
	frontendURL string
	auth        authConfig
}

type emailConfig struct {
	fromEmail string
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
	expiry    time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type mailTrapConfig struct {
	apiKey string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.putActivateUser)

			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.putFollowHandler)
				r.Put("/unfollow", app.putUnfollowHandler)

			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.postAuthenticateUserHandler)
			r.Post("/token", app.postTokenHandler)
		})
	})

	return r
}

func (a *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = a.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         a.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	a.logger.Infow("server started", "addr", a.config.addr, "env", a.config.env, "token exp", a.config.auth.jwt.exp)
	return srv.ListenAndServe()
}
