package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env            string `env:"ENV" envDefault:"local"`
	ServicePort    string `env:"SERVICE_PORT" envDefault:"8080"`
	SQLDatabaseURL string `env:"SQL_DATABASE_URL,required"`
	JWTSecret      string `env:"JWT_SECRET,file,required"`
	SmtpHost       string `env:"SMTP_HOST,required"`
	SmtpPort       string `env:"SMTP_PORT,required"`
	SmtpUser       string `env:"SMTP_USER,required"`
	SmtpPassword   string `env:"SMTP_PASSWORD,required"`
}

func Load() (*Config, error) {
	var (
		cfg Config
		err error
	)
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using global config")
	}
	if err = env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
