package main

import (
	"github.com/deriannavy/microgo/internal/store"
	"net/http"
)

type CreateTransactionPayload struct {
	AccountID int64  `json:"account_id"`
	Date      string `json:"date"`
	Amount    int32  `json:"amount"`
	// accountOut
	// accountIN
	Place       string   `json:"place"`
	Description string   `json:"description"`
	Tag         []string `json:"tag"`
}

func (app *application) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction store.Transaction
	if err := readJSON(w, r, transaction); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	accountID := 1

	ctx := r.Context()

	app.store.Transaction.Create(ctx)
}
