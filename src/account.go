package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/deriannavy/microgo/internal/store"
	"github.com/go-chi/chi/v5"
)

type accountKey string

const accountCtx accountKey = "account"

func (app *application) accountContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, accountCtx, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAccountFromCtx(r *http.Request) *store.Account {
	account, _ := r.Context().Value(accountCtx).(*store.Account)
	return account
}

// GetAccountHandler godoc
//
//	@Sumary			feteches
//	@Description	feteches
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account Id"
//	@Success		200	{object}	store.Account
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApyKeyAuth
//	@Router			/accounts/{id} [get]
func (app *application) getAccountHandler(w http.ResponseWriter, r *http.Request) {

	account := getAccountFromCtx(r)

	if err := writeJSON(w, http.StatusOK, account); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
