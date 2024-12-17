package main

import (
	"time"

	"github.com/edwrdc/digitally/internal/auth"
	"github.com/edwrdc/digitally/internal/db"
	"github.com/edwrdc/digitally/internal/env"
	"github.com/edwrdc/digitally/internal/mailer"
	"github.com/edwrdc/digitally/internal/store"
	"github.com/edwrdc/digitally/internal/store/cache"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Digitally API
//	@description	API for Digitally, a platform for buying and selling digital products.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				JWT authorization header

func main() {

	cfg := config{
		addr:        env.Get("API_PORT", ":8080"),
		env:         env.Get("API_ENV", "development"),
		apiURL:      env.Get("EXTERNAL_API_URL", "localhost:6969"),
		frontendURL: env.Get("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			dsn:          env.Get("DB_DSN", "postgres://admin:adminpassword@localhost:5432/digitally?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_IDLE_TIME", 15)) * time.Minute,
		},
		redisCfg: redisConfig{
			addr:    env.Get("REDIS_ADDR", "localhost:6379"),
			pw:      env.Get("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		mail: mailConfig{
			fromEmail: env.Get("MAIL_FROM_EMAIL", ""),
			exp:       time.Duration(env.GetInt("MAIL_EXPIRY", 3)) * time.Hour,
			mailtrap: mailtrapConfig{
				apiKey:  env.Get("MAILTRAP_API_KEY", ""),
				inboxID: env.Get("MAILTRAP_INBOX_ID", ""),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user: env.Get("BASIC_AUTH_USER", "admin"),
				pass: env.Get("BASIC_AUTH_PASS", "adminpassword"),
			},
			token: tokenAuthConfig{
				secret: env.Get("AUTH_JWT_SECRET", "digitallyio"),
				expiry: time.Duration(env.GetInt("AUTH_JWT_EXPIRY", 3)) * time.Hour * 24,
				iss:    "digitally",
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()

	logger.Info("Established connection pool to database")

	// Redis
	var redisDB *redis.Client
	if cfg.redisCfg.enabled {
		redisDB = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Infow("Redis cache connected")
	}

	store := store.New(db)
	cacheStorage := cache.NewRedisStorage(redisDB)

	mailer := mailer.NewMailtrapMailer(
		cfg.mail.mailtrap.apiKey,
		cfg.mail.fromEmail,
		cfg.mail.mailtrap.inboxID,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	app.logger.Infow("Server Started", "env", app.config.env, "addr", app.config.addr)

	logger.Fatal(app.run())

}
