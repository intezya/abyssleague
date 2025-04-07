package dotenv

import (
	"os"
	"strconv"
)

func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		if fallback == "" {
			panic("Missing environment variable " + key)
		}
		return fallback
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		if fallback == 0 {
			panic("Missing environment variable " + key)
		}
		return fallback
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return result
}
