package config

import (
	"context"
	"errors"
	"os"

	"github.com/abibby/salusa/env"
	"github.com/abibby/salusa/kernel"
	"github.com/joho/godotenv"
)

var Port int
var DBUsername string
var DBPassword string
var DBDatabase string
var DBHost string
var DBPort int
var EztvDomain string

func Load(ctx context.Context) error {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
	} else if err != nil {
		return err
	}

	Port = env.Int("PORT", 60575)

	DBUsername = env.String("DB_USERNAME", "eztv")
	DBPassword = env.String("DB_PASSWORD", "secret")
	DBDatabase = env.String("DB_DATABASE", "eztv")
	DBHost = env.String("DB_HOST", "127.0.0.1")
	DBPort = env.Int("DB_PORT", 3306)

	EztvDomain = env.String("EZTV_DOMAIN", "https://eztv.re")
	return nil
}
func Kernel() *kernel.KernelConfig {
	return &kernel.KernelConfig{
		Port: Port,
	}
}
