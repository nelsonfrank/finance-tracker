package main

import (
	"log"
	"time"

	"github.com/nelsonfrank/finance-tracker/internal/auth"
	"github.com/nelsonfrank/finance-tracker/internal/db"
	"github.com/nelsonfrank/finance-tracker/internal/env"
	"github.com/nelsonfrank/finance-tracker/internal/store"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "host=localhost user=admin password=adminpassword dbname=finance-tracker port=5438  sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		oAuth: oAuthConfig{
			google: &oauth2.Config{
				ClientID:     env.GetString("GOOGLE_CLIENT_ID", ""),
				ClientSecret: env.GetString("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  env.GetString("OAUTH_REDIRECT_URL", "http://localhost:3000/v1/auth/google/callback"),
				Scopes:       []string{"email", "profile"},
				Endpoint:     google.Endpoint,
			},
		},
		mfa: mfaConfig{
			token: jwtToken{
				secret:          env.GetString("JWT_SECRET", ""),
				refreshTokenExp: time.Hour * 24 * 3,
				exp:             time.Second * 5,
				iss:             "financial-tracker"},
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("database connection pool established")

	store := store.NewStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.mfa.token.secret,
		cfg.mfa.token.iss,
		cfg.mfa.token.iss,
	)
	app := &application{
		config:        cfg,
		store:         store,
		db:            db,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	log.Fatal((app.run(mux)))
}
