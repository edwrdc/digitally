package env

import (
	"os"
	"strconv"
)

func Get(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func GetInt(key string, fallback int) int {
	v := Get(key, "")
	if v == "" {
		return fallback
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
