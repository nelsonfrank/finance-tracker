package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var randomState string = "random"

// ValidationError represents a custom error response
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
type LoginUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		validationErrors := app.validationErrorFormatter(err)

		sendError(w, http.StatusBadRequest, validationErrors)
		return
	}

	log.Print(payload)

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		validationErrors := app.validationErrorFormatter(err)

		sendError(w, http.StatusBadRequest, validationErrors)
		return
	}

	log.Print(payload)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

}

// google OAuth2
func (app *application) oAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := app.config.oAuth.google.AuthCodeURL(randomState, oauth2.AccessTypeOffline)

	w.Write([]byte(url))
	// http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// Exchanging the code for an access token
	t, err := app.config.oAuth.google.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Creating an HTTP client to make authenticated request using the access key.
	// This client method also regenerate the access key using the refresh key.
	client := app.config.oAuth.google.Client(context.Background(), t)

	// Getting the user public details from google API endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Closing the request body when this function returns.
	// This is a good practice to avoid memory leak
	defer resp.Body.Close()

	var v any

	// Reading the JSON body using JSON decoder
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// sending the user public value as a response. This is may not be a good practice,
	// but for demonstration, I think it serves the need.
	fmt.Fprintf(w, "%v", v)
}
