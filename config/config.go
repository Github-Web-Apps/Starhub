package config

import (
	"github.com/apex/log"
	"github.com/caarlos0/env"
)

// Config of the app
type Config struct {
	Port           string `env:"PORT" envDefault:"3000"`
	ClientID       string `env:"GITHUB_CLIENT_ID"`
	ClientSecret   string `env:"GITHUB_CLIENT_SECRET"`
	OauthState     string `env:"OAUTH_STATE"`
	DatabaseURL    string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/watchub?sslmode=disable"`
	SendgridAPIKey string `env:"SENDGRID_API_KEY"`
	Schedule       string `env:"SCHEDULE" envDefault:"@every 1m"`
	SessionSecret  string `env:"SESSION_SECRET" envDefault:"super-secret-session-secret"`
	SessionName    string `env:"SESSION_NAME" envDefault:"JSESSIONID"`
}

// Get the config
func Get() (cfg Config) {
	if err := env.Parse(&cfg); err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
	return
}
