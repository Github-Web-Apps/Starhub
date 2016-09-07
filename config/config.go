package config

import "github.com/caarlos0/env"

// Config of the app
type Config struct {
	Port           int    `env:"PORT" envDefault:"3000"`
	ClientID       string `env:"GITHUB_CLIENT_ID"`
	ClientSecret   string `env:"GITHUB_CLIENT_SECRET"`
	OauthState     string `env:"OAUTH_STATE"`
	DatabaseURL    string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/watchub?sslmode=disable"`
	SendgridAPIKey string `env:"SENDGRID_API_KEY"`
	Schedule       string `env:"SCHEDULE" envDefault:"@hourly"`
}

// Get the config
func Get() (Config, error) {
	cfg := Config{}
	return cfg, env.Parse(&cfg)
}
