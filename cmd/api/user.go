package main

import (
	"net/http"

	"github.com/nelsonfrank/finance-tracker/internal/store"
)

type userKey string

const userCtx userKey = "user"

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
