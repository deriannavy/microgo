package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/deriannavy/microgo/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) checkRolePrecedence(ctx context.Context, account *store.Account, roleName string) (bool, error) {
	role, err := app.store.Role.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return account.Role.Level >= role.Level, nil
}

func (app *application) checkTransactionOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		account := getAccountFromCtx(r)
		transaction := getTransactionFromCtx(r)

		if transaction.AccountId == account.Id {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), account, role)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenErrorResponse(w, r, fmt.Errorf("You are not allowed to perform this action"))
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("Autorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("Autorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("Autorization header is malformed"))
				return
			}

			username := app.config.auth.basic.user
			password := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("Invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedErrorResponse(w, r, fmt.Errorf("Autorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("Autorization header is malformed"))
				return
			}

			token := parts[1]
			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			claims := jwtToken.Claims.(jwt.MapClaims)

			accountId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx := r.Context()

			account, err := app.store.Account.GetById(ctx, accountId)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, accountCtx, account)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
