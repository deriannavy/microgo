package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error %s Path: %s Error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusNotFound, "Not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized basic Error %s Path: %s Error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) forbiddenErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized basic Error %s Path: %s Error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusForbidden, "Forbidden")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized basic Error %s Path: %s Error: %s", r.Method, r.URL.Path, err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) rateLimitExceededReponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	log.Printf("Rate limit exceed Path: %s Error: %s", r.Method, r.URL.Path)

	w.Header().Set("X-Retry-After", retryAfter)

	writeJSON(w, http.StatusTooManyRequests, "Rate limit Exceed, retry after"+retryAfter)
}
