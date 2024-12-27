package main

import (
	"net/http"
)

type api struct {
	addr string
}

func (s *api) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Account List"))
}

func (s *api) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Account Created"))
}

func main() {
	api := &api{addr: ":8080"}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /accounts", api.getAccountHandler)
	mux.HandleFunc("POST /accounts", api.createAccountHandler)

	err := srv.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
