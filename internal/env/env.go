package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetDuration(key, fallback string) time.Duration {
	val, ok := os.LookupEnv(key)

	if !ok {
		fallbackDur, err := time.ParseDuration(fallback)
		if err != nil {
			return 0
		}
		return fallbackDur
	}

	actualDuration, err := time.ParseDuration(val)

	if err != nil {
		return 0
	}

	return actualDuration
}
