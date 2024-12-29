package main

import (
	"time"

	"github.com/nelsonfrank/finance-tracker/internal/auth"
	"github.com/nelsonfrank/finance-tracker/internal/db"
	"github.com/nelsonfrank/finance-tracker/internal/env"
	"github.com/nelsonfrank/finance-tracker/internal/mailer"
	"github.com/nelsonfrank/finance-tracker/internal/store"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {

	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		env:         env.GetString("ENV", "development"),
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
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.mfa.token.secret,
		cfg.mfa.token.iss,
		cfg.mfa.token.iss,
	)

	// Mailer
	// mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config:        cfg,
		store:         store,
		db:            db,
		authenticator: jwtAuthenticator,
		mailer:        mailtrap,
		logger:        logger,
	}

	mux := app.mount()

	logger.Fatal((app.run(mux)))
}
