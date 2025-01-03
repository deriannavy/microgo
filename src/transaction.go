package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/deriannavy/microgo/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateTransactionPayload struct {
	AccountId int64  `json:"account_id"`
	Date      string `json:"date"`
	Amount    int32  `json:"amount"`
	// accountOut
	// accountIN
	Place       string   `json:"place" validate:"required,max=100"`
	Description string   `json:"description" validate:"required,max=100"`
	Tag         []string `json:"tag"`
}

type transactionKey string

const transactionCtx transactionKey = "transaction"

func (app *application) transactionContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "transactionId")
		transactionId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		transaction, err := app.store.Transaction.GetById(ctx, transactionId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, transactionCtx, transaction)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTransactionFromCtx(r *http.Request) *store.Transaction {
	transaction, _ := r.Context().Value(transactionCtx).(*store.Transaction)
	return transaction
}

func (app *application) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateTransactionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := &store.Transaction{
		AccountId:   1,
		Date:        payload.Date,
		Amount:      payload.Amount,
		Place:       payload.Place,
		Description: payload.Description,
		Tag:         payload.Tag,
	}

	ctx := r.Context()

	if err := app.store.Transaction.Create(ctx, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getTransactionHandler(w http.ResponseWriter, r *http.Request) {

	transaction := getTransactionFromCtx(r)

	if err := writeJSON(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type UpdateTransactionPayload struct {
	Date   *string `json:"date"`
	Amount *int32  `json:"amount"`
	// accountOut
	// accountIN
	Place       *string   `json:"place" validate:"omitempty,max=100"`
	Description *string   `json:"description" validate:"omitempty,max=100"`
	Tag         *[]string `json:"tag"`
}

func (app *application) patchTransactionHandler(w http.ResponseWriter, r *http.Request) {

	transaction := getTransactionFromCtx(r)
	var payload UpdateTransactionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Date != nil {
		transaction.Date = *payload.Date
	}
	if payload.Amount != nil {
		transaction.Amount = *payload.Amount
	}
	if payload.Place != nil {
		transaction.Place = *payload.Place
	}
	if payload.Description != nil {
		transaction.Description = *payload.Description
	}
	if payload.Tag != nil {
		transaction.Tag = *payload.Tag
	}

	if err := app.store.Transaction.Update(r.Context(), transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "transactionId")
	transactionId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Transaction.Delete(ctx, transactionId); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
