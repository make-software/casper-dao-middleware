package config

import (
	"log"
	"os"
)

type DBConfig struct {
	DatabaseURI             string `env:"DATABASE_URI,required"`
	MaxAllowedDBBufferBytes uint64 `env:"MAX_ALLOWED_DB_BUFFER_BYTES" envDefault:"0"`
	MaxOpenConnections      int    `env:"DATABASE_MAX_OPEN_CONNECTIONS" envDefault:"5"`
	MaxIdleConnections      int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" envDefault:"5"`
}

func GetEnv(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Panicf("missing required ENV: %s", envName)
	}
	return env
}
