package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Load() (Config, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetDefault("runtime.env", "local")

	v.SetDefault("http.port", 8080)
	v.SetDefault("http.read_header_timeout", "5s")
	v.SetDefault("http.read_timeout", "5s")
	v.SetDefault("http.write_timeout", "10s")
	v.SetDefault("http.idle_timeout", "120s")
	v.SetDefault("http.shutdown_timeout", "10s")
	v.SetDefault("http.max_header_bytes", 1048576)

	v.SetDefault("db.url", "postgres://blog:blog@postgres:5432/blog?sslmode=disable")
	v.SetDefault("db.max_conns", 10)
	v.SetDefault("db.min_conns", 2)
	v.SetDefault("db.max_conn_lifetime", "1h")
	v.SetDefault("db.max_conn_idle_time", "30m")

	v.SetDefault("jwt.secret", "dev-only-secret-change-me-in-production-please-32b")
	v.SetDefault("jwt.ttl", "24h")
	v.SetDefault("jwt.issuer", "blog")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	v.SetEnvPrefix("BLOG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func validate(c Config) error {
	switch c.Runtime.Env {
	case "local", "staging", "production":
	default:
		return fmt.Errorf("runtime.env must be local|staging|production")
	}
	if c.DB.URL == "" {
		return fmt.Errorf("db.url is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("jwt.secret must be at least 32 bytes")
	}
	if c.Runtime.Env == "production" && strings.Contains(c.JWT.Secret, "dev-only") {
		return fmt.Errorf("jwt.secret must not use the dev default in production")
	}
	if c.HTTP.MaxHeaderBytes <= 0 {
		return fmt.Errorf("http.max_header_bytes must be positive")
	}
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be debug|info|warn|error")
	}
	return nil
}
