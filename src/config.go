package main

import (
	"github.com/deriannavy/microgo/internal/auth"
	"github.com/deriannavy/microgo/internal/cache"
	"github.com/deriannavy/microgo/internal/rateLimiter"
	"github.com/deriannavy/microgo/internal/store"
	"time"
)

type application struct {
	config config
	store  store.Storage
	cache  cache.Storage
	//mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   rateLimiter.RateLimiter
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	exp       time.Duration
	fromEmail string
}

type sendGridConfig struct {
	apiKey string
}

type config struct {
	addr          string
	env           string
	version       string
	allowedOrigin string
	apiURL        string
	frontURL      string
	apiVersion    string
	db            dbConfig
	mailer        mailConfig
	auth          authConfig
	cache         cacheConfig
	rateLimiter   rateLimiterConfig
}

type cacheConfig struct {
	addr    string
	pass    string
	db      int
	enabled bool
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type rateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
