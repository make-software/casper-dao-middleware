package boot

import (
	"log"

	"github.com/joho/godotenv"
)

type Config interface {
	Parse() error
}

const EnvFile = ".env"

func ParseEnvConfig(config Config) error {
	if err := godotenv.Load(EnvFile); err != nil {
		log.Printf("cant load from .env %s, trying load from ENV\n", err.Error())
	}

	return config.Parse()
}
