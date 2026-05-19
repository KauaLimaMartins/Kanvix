package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr          string
	DBPath        string
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	SessionTTL time.Duration
	CacheTTL   time.Duration

	CookieName   string
	CookieDomain string
	CookieSecure bool

	AllowedOrigins []string
}

func Load() Config {
	return Config{
		Addr:           getenv("KANVIX_ADDR", ":8080"),
		DBPath:         getenv("KANVIX_DB_PATH", "./data/kanvix.db"),
		RedisAddr:      getenv("KANVIX_REDIS_ADDR", "redis:6379"),
		RedisPassword:  getenv("KANVIX_REDIS_PASSWORD", ""),
		RedisDB:        getenvInt("KANVIX_REDIS_DB", 0),
		SessionTTL:     getenvDuration("KANVIX_SESSION_TTL", 7*24*time.Hour),
		CacheTTL:       getenvDuration("KANVIX_CACHE_TTL", 30*time.Second),
		CookieName:     getenv("KANVIX_COOKIE_NAME", "kanvix_session"),
		CookieDomain:   getenv("KANVIX_COOKIE_DOMAIN", ""),
		CookieSecure:   getenvBool("KANVIX_COOKIE_SECURE", false),
		AllowedOrigins: getenvCSV("KANVIX_ALLOWED_ORIGINS", ""),
	}
}

func getenv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getenvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getenvBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func getenvDuration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func getenvCSV(key, def string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		raw = def
	}
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
