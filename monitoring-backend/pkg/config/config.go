package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	ApiKey  string

	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	DBSSL    string
}

func Load() Config {
	_ = godotenv.Load() // ไม่ error แม้ไม่มีไฟล์
	cfg := Config{
		AppPort: env("APP_PORT", "8080"),
		ApiKey:  env("API_KEY", "supersecretapikey"),
		DBHost:  env("DB_HOST", "localhost"),
		DBPort:  env("DB_PORT", "5432"),
		DBUser:  env("DB_USER", "postgres"),
		DBPass:  env("DB_PASSWORD", "postgres"),
		DBName:  env("DB_NAME", "metrics"),
		DBSSL:   env("DB_SSLMODE", "disable"),
	}
	if cfg.ApiKey == "" { log.Println("[WARN] API_KEY is empty") }
	return cfg
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}