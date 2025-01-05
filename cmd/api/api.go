package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nelsonfrank/finance-tracker/internal/auth"
	"github.com/nelsonfrank/finance-tracker/internal/mailer"
	"github.com/nelsonfrank/finance-tracker/internal/store"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type application struct {
	config        config
	store         store.Storage
	db            *gorm.DB
	authenticator auth.Authenticator
	mailer        mailer.Client
	logger        *zap.SugaredLogger
}

type config struct {
	addr        string
	frontendURL string
	env         string
	db          dbConfig
	oAuth       oAuthConfig
	mfa         mfaConfig
	mail        mailConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type oAuthConfig struct {
	google *oauth2.Config
}

type mfaConfig struct {
	token jwtToken
}

type jwtToken struct {
	secret          string
	iss             string
	exp             time.Duration
	refreshTokenExp time.Duration
}

type mailConfig struct {
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
	fromEmail string
	exp       time.Duration
}

type mailTrapConfig struct {
	apiKey string
}

type sendGridConfig struct {
	apiKey string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/auth", func(r chi.Router) {
			// google OAuth2
			r.Get("/google", app.oAuthHandler)
			r.Get("/google/callback", app.oAuthCallbackHandler)

			// MFA
			r.Post("/register", app.register)
			r.Post("/login", app.login)
			r.Post("/logout", app.logout)
			r.Group(func(r chi.Router) {
				r.Use(app.RefreshTokenMiddleware)
				r.Post("/refresh-token", app.refreshTokenHandler)
			})
		})

		r.Route("/dashboard", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/", app.dashboardHandler)
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)
	return srv.ListenAndServe()
}
