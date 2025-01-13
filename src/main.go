package main

import (
	"github.com/deriannavy/microgo/internal/cache"
	"github.com/deriannavy/microgo/internal/rateLimiter"
	"github.com/go-redis/redis/v8"
	"log"
	"time"

	"github.com/deriannavy/microgo/internal/auth"
	"github.com/deriannavy/microgo/internal/db"
	"github.com/deriannavy/microgo/internal/env"

	// "github.com/deriannavy/microgo/internal/mailer"
	"github.com/deriannavy/microgo/internal/store"
)

//	@title			Swagger Example API
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApikeyAuth
// @in							header
// @name						Autorization
// @description
func main() {
	cfg := config{
		addr:          env.GetEnvString("ADDR", ":8080"),
		env:           env.GetEnvString("ENV", "development"),
		version:       env.GetEnvString("VERSION", "1.0.0"),
		apiURL:        env.GetEnvString("API_URL", "localhost:8080"),
		frontURL:      env.GetEnvString("API_URL", "http://localhost:8080"),
		apiVersion:    env.GetEnvString("API_VERSION", "/v1"),
		allowedOrigin: env.GetEnvString("CORS_ALLOWED_ORIGIN", "http://localhost:5174"),
		db: dbConfig{
			addr:         env.GetEnvString("DB_ADDR", "postgresql://accounts:4cc0unts@localhost:5430/finance?sslmode=disable"),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetEnvString("DB_MAX_IDLE_TIME", "15m"),
		},
		mailer: mailConfig{
			fromEmail: env.GetEnvString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetEnvString("API_KEY", ""),
			},
			exp: time.Hour * 24 * 3, // 3 days
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetEnvString("AUTH_BASIC_USER", "admin"),
				pass: env.GetEnvString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetEnvString("AUTH_TOKEN_SECRET", "admin"),
				exp:    time.Hour * 24 * 1, // 1 day
				iss:    env.GetEnvString("AUTH_TOKEN_ISS", "finance"),
			},
		},
		cache: cacheConfig{
			addr:    env.GetEnvString("CACHE_ADDR", "localhost:6379"),
			pass:    env.GetEnvString("CACHE_PASS", ""),
			db:      env.GetEnvInt("CACHE_DB", 0),
			enabled: env.GetEnvBool("CACHE_ENABLED", true),
		},
		rateLimiter: rateLimiterConfig{
			RequestsPerTimeFrame: env.GetEnvInt("RATE_LIMITER_REQUEST_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetEnvBool("RATE_LIMITER_ENABLE", true),
		},
	}

	newDB, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}
	defer newDB.Close()
	log.Println("Database connected...")

	storage := store.NewStorage(newDB)

	var cdb *redis.Client
	if cfg.cache.enabled {
		cdb = db.NewCacheClient(cfg.cache.addr, cfg.cache.pass, cfg.cache.db)
		log.Println("Cache database connected...")
	}

	cacheStorage := cache.NewCacheStorage(cdb)

	//mail := mailer.NewSendGrid(
	//	cfg.mailer.sendGrid.apiKey,
	//	cfg.mailer.fromEmail,
	//)

	rl := rateLimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config: cfg,
		store:  storage,
		cache:  cacheStorage,
		//mailer: mail,
		authenticator: jwtAuthenticator,
		rateLimiter:   rl,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
