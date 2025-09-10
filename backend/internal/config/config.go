package config

import (
	"fmt"
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
		Enabled  bool
		URL      string
		Password string
		DB       int
	}

	Security struct {
		JWTSecret string
		APIKey    string
	}
}

// Load reads configuration from .env, config.yaml, env variables, and defaults
func Load() (*Config, error) {
	// ----------------------------
	// Load .env if exists
	// ----------------------------
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables or defaults")
	}

	// ----------------------------
	// Setup Viper
	// ----------------------------
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config") // support subfolder

	// Automatic environment variable override
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
		log.Printf("Config file not found, using env/defaults: %v", err)
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// ----------------------------
	// Unmarshal into struct
	// ----------------------------
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// ----------------------------
	// Optional validation
	// ----------------------------
	if cfg.App.Port == "" {
		cfg.App.Port = "8080"
	}

	log.Printf("App config loaded: Env=%s, Port=%s, DB=%s, Redis Enabled=%v",
		cfg.App.Env, cfg.App.Port, cfg.Database.URL, cfg.Redis.Enabled)

	return &cfg, nil
}
