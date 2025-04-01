package env

import (
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}
	return env
}

func GetInt(key string, fallback int) int {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}

	valAsInt, err := strconv.Atoi(env)

	if err != nil {
		return fallback
	}

	return valAsInt
}
func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)

	if err != nil {
		return fallback
	}

	return boolVal
}
