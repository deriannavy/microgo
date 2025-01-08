package main

import (
	"log"
	"time"

	"github.com/deriannavy/microgo/internal/db"
	"github.com/deriannavy/microgo/internal/env"
	// "github.com/deriannavy/microgo/internal/mailer"
	"github.com/deriannavy/microgo/internal/store"
)

const version = "0.0.1"

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
		addr: env.GetEnvString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetEnvString("DB_ADDR", "postgresql://accounts:4cc0unts@localhost:5430/finance?sslmode=disable"),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetEnvString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:        env.GetEnvString("ENV", "development"),
		apiURL:     env.GetEnvString("API_URL", "localhost:8080"),
		frontURL:   env.GetEnvString("API_URL", "http://localhost:8080"),
		apiVersion: env.GetEnvString("API_VERSION", "/v1"),
		mailer: mailConfig{
			fromEmail: env.GetEnvString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetEnvString("API_KEY", ""),
			},
			exp: time.Hour * 24 * 3, // 3 days
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

	//mail := mailer.NewSendGrid(
	//	cfg.mailer.sendGrid.apiKey,
	//	cfg.mailer.fromEmail,
	//)

	app := &application{
		config: cfg,
		store:  storage,
		//mailer: mail,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
