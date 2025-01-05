package main

import (
	"net/http"

	"github.com/nelsonfrank/finance-tracker/internal/db/model"
)

type userKey string

const userCtx userKey = "user"

func getUserFromContext(r *http.Request) model.User {
	user, _ := r.Context().Value(userCtx).(model.User)
	return user
}
