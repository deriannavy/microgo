package main

import (
	"github.com/deriannavy/microgo/internal/db"
	"github.com/deriannavy/microgo/internal/env"
	"github.com/deriannavy/microgo/internal/store"
	"log"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetEnvString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetEnvString("DB_ADDR", "postgresql://accounts:4cc0unts@localhost:5430/finance?sslmode=disable"),
			maxOpenConns: env.GetEnvInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetEnvInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetEnvString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetEnvString("ENV", "development"),
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

	app := &application{
		config: cfg,
		store:  storage,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
