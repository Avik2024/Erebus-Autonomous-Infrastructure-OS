package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App struct {
		Name string
		Env  string
		Port string
	}

	Database struct {
		URL         string
		Migrations  string
		AutoMigrate bool
	}

	Redis struct {
		URL      string
		Password string
		DB       int
		Enabled  bool
	}

	Security struct {
		JWTSecret string
		APIKey    string
	}
}

// Load reads configuration from .env, config.yaml, env variables, and defaults
func Load() (*Config, error) {
	// Load .env file if present
	_ = godotenv.Load()

	// Setup Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config") // support subfolder

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ----------------------------
	// Defaults
	// ----------------------------
	viper.SetDefault("app.name", "erebus")
	viper.SetDefault("app.env", "dev")
	viper.SetDefault("app.port", "8080")

	viper.SetDefault("database.url", "postgres://erebus:erebus@postgres:5432/erebus_db?sslmode=disable")
	viper.SetDefault("database.migrations", "./migrations")
	viper.SetDefault("database.automigrate", false)

	viper.SetDefault("redis.enabled", true)
	viper.SetDefault("redis.url", "redis:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("security.jwtsecret", "changeme")
	viper.SetDefault("security.apikey", "")

	// ----------------------------
	// Read config.yaml if exists
	// ----------------------------
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file not found, using env/defaults: %v", err)
	}

	// ----------------------------
	// Unmarshal into struct
	// ----------------------------
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
