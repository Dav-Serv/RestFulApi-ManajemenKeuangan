package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port				string
	GinMode				string
	DBPath				string
	JWTSecret			string
	JWTExpiryHrs		string
	GoogleClientID		string
	GoogleClientSecret	string
	GoogleRedirectURL	string
	FrontendURL			string
}

var App *Config

// Load membaca file .env (jika ada) lalu mengisi struct Config.
// Dipanggil sekali di awal main().
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env tidak ditemukan, menggunakan enviroment variable sistem")
	}

	App = &Config{
		Port:						getEnv("PORT", "8000"),
		GinMode: 					getEnv("GIN_MODE", "debug"),
		DBPath: 					getEnv("DB_PATH", "expense.db"),
		JWTSecret: 					getEnv("JWT_SECRET", "secret-default-jangan-dipakai-di-production"),
		JWTExpiryHrs: 				getEnv("JWT_EXPIRY_HOURS", "72"),
		GoogleClientID: 			getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: 		getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL: 			getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		FrontendURL: 				getEnv("FRONTEND_URL", "http://localhost:3000/auth/callback"),
	}

	return App
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}