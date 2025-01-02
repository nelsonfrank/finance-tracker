package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nelsonfrank/finance-tracker/internal/db/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var randomState string = "random"

// ValidationError represents a custom error response
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=3,max=72"`
}
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type LoginResponse struct {
	User model.User `json:"user"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		validationErrors := app.validationErrorFormatter(err)

		sendError(w, http.StatusBadRequest, validationErrors)
		return
	}

	// Check if user already exists
	var existingUser model.User
	result := app.db.Where("email = ?", payload.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	// Create new user
	user := model.User{
		Email:     payload.Email,
		Password:  string(hashedPassword),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}

	if result := app.db.Create(&user); result.Error != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, user)

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		validationErrors := app.validationErrorFormatter(err)

		sendError(w, http.StatusBadRequest, validationErrors)
		return
	}

	var user model.User
	result := app.db.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			writeJSONError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		app.internalServerError(w, r, result.Error)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT Access token
	claims := app.authenticator.JwtClaimGenerator(
		user.ID,
		app.config.mfa.token.exp,
		app.config.mfa.token.iss,
		app.config.mfa.token.iss,
	)

	accessToken, err := app.authenticator.GenerateToken(claims)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Generate JWT Refresh token
	refreshTokenClaims := app.authenticator.JwtClaimGenerator(
		user.ID,
		app.config.mfa.token.refreshTokenExp,
		app.config.mfa.token.iss,
		app.config.mfa.token.iss)

	refreshToken, err := app.authenticator.GenerateToken(refreshTokenClaims)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error generating token")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(app.config.mfa.token.exp),
		HttpOnly: false,
		Secure:   false, // Set to true in production (requires HTTPS)
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(app.config.mfa.token.refreshTokenExp),
		HttpOnly: true,
		Secure:   false, // Set to true in production (requires HTTPS)
		SameSite: http.SameSiteNoneMode,
	})
	writeJSON(w, http.StatusOK, &LoginResponse{
		user,
	})

}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

}

func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload RefreshTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		validationErrors := app.validationErrorFormatter(err)

		sendError(w, http.StatusBadRequest, validationErrors)
		return
	}

	jwtToken, err := app.authenticator.ValidateToken(payload.RefreshToken)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	claims, _ := jwtToken.Claims.(jwt.MapClaims)

	userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// Generate JWT Access token
	newClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(app.config.mfa.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.mfa.token.iss,
		"aud": app.config.mfa.token.iss,
	}
	accessToken, err := app.authenticator.GenerateToken(newClaims)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	writeJSON(w, http.StatusOK, &RefreshTokenResponse{
		accessToken,
	})
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
