package main

import (
	"github.com/deriannavy/microgo/internal"
	"log"
)

func main() {
	cfg := config{
		addr: env.GetEnvString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
