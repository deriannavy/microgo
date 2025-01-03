package main

import (
	"github.com/deriannavy/microgo/internal/store"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "accountId")
	accountId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	account, err := app.store.Account.GetById(ctx, accountId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:

			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, account); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
