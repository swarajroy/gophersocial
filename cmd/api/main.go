package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/swarajroy/gophersocial/internal/db"
	"github.com/swarajroy/gophersocial/internal/env"
	"github.com/swarajroy/gophersocial/internal/mailer"
	"github.com/swarajroy/gophersocial/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for Gopher social network.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	    /v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		env:         env.GetString("ENV", "dev"),
		dbConfig: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", "15m"),
		},
		email: emailConfig{
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
			expiry: env.GetDuration("EMAIL_INVITATION_EXPIRY", "24h"),
		},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	db, err := db.New(cfg.dbConfig.addr, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxIdleTime)
	if err != nil {
		logger.Fatal("connection to db failed!")
	}
	defer db.Close()

	logger.Info("database connections pool configured successfully!")

	// store
	store := store.NewStorage(db)

	// mailer
	templateBuilder := mailer.NewTemplateBuilder()
	//mailer, err := mailer.NewSendGridMailer(cfg.email.fromEmail, cfg.email.sendGrid.apiKey, templateBuilder)
	mailer, err := mailer.NewMailtrapMailer(cfg.email.fromEmail, cfg.email.mailTrap.apiKey, templateBuilder)
	if err != nil {
		logger.Fatal("mailer configuration failed!")
	}

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
