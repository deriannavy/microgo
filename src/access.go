package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/deriannavy/microgo/internal/mailer"
	"github.com/deriannavy/microgo/internal/store"
	"github.com/google/uuid"
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
	if err := readJSON(w, r, &payload); err != nil {
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

	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Account.CreateAndConfirm(ctx, account, hashToken, app.config.mailer.exp)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	activationURL := fmt.Stringf("%s/confirm/%s", app.config.frontURL, hashToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      account.Username,
		ActivationURL: activationURL,
	}

	// send mail
	err = app.mailer.Send(
		mailer.AccountWelcomeTemplate,
		account.Username,
		account.Email,
		vars,
		false,
	)

	if err != nil {

		if err := app.store.Account.Delete(ctx, account.Id); err != nil {
			app.internalServerError(w, r, err)
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
