package config

import "os"

type Config struct {
    Port     string
    Env      string
    DBHost   string
    DBPort   string
    DBUser   string
    DBPass   string
    DBName   string
	JWTSecret string
}

func Load() *Config {
    return &Config{
        Port:   getEnv("PORT", "8080"),
        Env:    getEnv("ENV", "development"),
        DBHost: getEnv("DB_HOST", "localhost"),
        DBPort: getEnv("DB_PORT", "5432"),
        DBUser: getEnv("DB_USER", "postgres"),
        DBPass: getEnv("DB_PASSWORD", "password"),
        DBName:   getEnv("DB_NAME", "mydb"),
        JWTSecret: getEnv("JWT_SECRET", "your-super-secret-key"),
    }
}

func getEnv(key, fallback string) string {
    if val, exists := os.LookupEnv(key); exists {
        return val
    }
    return fallback
}
