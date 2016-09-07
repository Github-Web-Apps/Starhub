package config

import "github.com/caarlos0/env"

type Config struct {
	Port         int    `env:"PORT" envDefault:"3000"`
	ClientID     string `env:"GITHUB_CLIENT_ID"`
	ClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	OauthState   string `env:"OAUTH_STATE"`
	DatabaseURL  string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/watchub?sslmode=disable"`
}

func Get() (Config, error) {
	cfg := Config{}
	return cfg, env.Parse(&cfg)
}
