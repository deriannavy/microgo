package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/deriannavy/microgo/docs"
	"github.com/deriannavy/microgo/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config config
	store  store.Storage
}

type mailConfig struct {
	exp time.Duration
}

type config struct {
	addr       string
	db         dbConfig
	env        string
	apiURL     string
	apiVersion string
	mail       mailConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped
	r.Use(middleware.Timeout(60 * time.Second))

	// D O C U M E N T A T I O N
	docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

	// H E A L T H
	r.Get("/health", app.healthCheckHandler)
	// V 1   R O U T E R
	r.Route(app.config.apiVersion, func(r chi.Router) {
		// A C C E S S   R O U T E R  --  P U B L I C
		r.Post("/register", app.registerAccountHandler)
		// A C C O U N T   R O U T E R
		r.Route("/account", func(r chi.Router) {
			// A C C O U N T   A C T I V A T E
			r.Put("/activate/{token}", app.activateAccountHandler)

			r.Route("/{accountId}", func(r chi.Router) {
				// M I D D L E W A R E   A C C O U N T
				r.Use(app.accountContextMiddleware)
				// G E T   A C C O U N T
				r.Get("/", app.getAccountHandler)
			})
		})

		// T R A N S A C T I O N   R O U T E R
		r.Route("/transaction", func(r chi.Router) {
			// G E T   I N D E X   T R A N S A C T I O N
			r.Post("/", app.getIndexTransactionHandler)
			// P O S T   T R A N S A C T I O N
			r.Post("/", app.createTransactionHandler)
			r.Route("/{transactionId}", func(r chi.Router) {
				// M I D D L E W A R E   T R A N S A C T I O N
				r.Use(app.transactionContextMiddleware)
				// G E T   T R A N S A C T I O N
				r.Get("/", app.getTransactionHandler)
				// P A T C H   T R A N S A C T I O N
				r.Patch("/", app.patchTransactionHandler)
				// D E L E T E   T R A N S A C T I O N
				r.Delete("/", app.deleteTransactionHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = app.config.apiVersion

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server at %s", app.config.addr)

	return srv.ListenAndServe()
}
