package main

import (
	"github.com/deriannavy/microgo/internal/store"
	"net/http"
)

type AccessPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// registerAccount godoc
//
// @Sumary asas
// @Description asx
// @Accept json
// @Produce json
// @Param payload body AccessPayload true "Account Credentials"
// @Success 201 {object} store.Account "Account Registered"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /register
func (app *application) registerAccountHandler(w http.ResponseWriter, r *http.Request) {
	var payload AccessPayload
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	account := &store.Account{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := account.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Account.CreateAndConfirm(ctx, account, "asa"); err != nil {

	}

	if err := writeJSON(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
