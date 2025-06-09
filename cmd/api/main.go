package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/swarajroy/gophersocial/internal/auth"
	"github.com/swarajroy/gophersocial/internal/db"
	"github.com/swarajroy/gophersocial/internal/env"
	"github.com/swarajroy/gophersocial/internal/mailer"
	"github.com/swarajroy/gophersocial/internal/store"
	"github.com/swarajroy/gophersocial/internal/store/cache"
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
		cache: cacheConfig{
			redisCfg: redisConfig{
				addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
				pw:      env.GetString("REDIS_PASSWORD", ""),
				db:      env.GetInt("REDIS_DB", 0),
				enabled: env.GetBool("REDIS_ENABLED", false),
			},
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
		auth: authConfig{
			basic: basicAuthConfig{
				user: env.GetString("BASIC_AUTH_USER", "admin"),
				pass: env.GetString("BASIC_AUTH_PASS", "admin"),
			},
			jwt: jwtConfig{
				secret: env.GetString("JWT_AUTH_SECRET", "example"), //TODO - dont ship defaults to PROD
				host:   env.GetString("JWT_HOST", "gophersocial"),
				exp:    env.GetDuration("JWT_AUTH_TOKEN_EXP", "24h"),
			},
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

	// caching
	var cacheStorage cache.Storage
	if cfg.cache.redisCfg.enabled {
		redisClient := cache.NewRedisClient(cfg.cache.redisCfg.addr, cfg.cache.redisCfg.pw, cfg.cache.redisCfg.db)
		cacheStorage = cache.NewRedisStorage(redisClient)
		logger.Infow("redis cache storage connection successul!")
		defer redisClient.Close()
	}

	// mailer
	templateBuilder := mailer.NewTemplateBuilder()
	//mailer, err := mailer.NewSendGridMailer(cfg.email.fromEmail, cfg.email.sendGrid.apiKey, templateBuilder)
	mailer, err := mailer.NewMailtrapMailer(cfg.email.fromEmail, cfg.email.mailTrap.apiKey, templateBuilder)
	if err != nil {
		logger.Fatal("mailer configuration failed!")
	}

	// stateless jwt token generator
	tokenGenerator := auth.NewJWTTokenGenerator(cfg.auth.jwt.secret, cfg.auth.jwt.host, cfg.auth.jwt.exp)
	app := &application{
		config: cfg,
		store:  store,
		cache:  cacheStorage,
		logger: logger,
		mailer: mailer,
		auth:   tokenGenerator,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
