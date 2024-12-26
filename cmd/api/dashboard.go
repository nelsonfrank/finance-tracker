package main

import (
	"net/http"
)

func (app *application) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, "hello there")
}
