package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type api struct {
	addr string
}

var accounts = []Account{}

func (s *api) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *api) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload Account
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a := Account{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}

	if err = insertAccount(a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func insertAccount(a Account) error {
	if a.FirstName == "" {
		return errors.New("First name is required")
	}

	if a.LastName == "" {
		return errors.New("Last name is required")
	}

	for _, account := range accounts {
		if account.FirstName == a.FirstName && account.LastName == a.LastName {
			return errors.New("Account already exists")
		}
	}

	accounts = append(accounts, a)
	return nil
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
