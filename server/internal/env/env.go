package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetOrDefault(key string, value string) string {
	target, found := os.LookupEnv(key)
	if !found {
		return value
	}
	return target
}

func GetOrDefaultInt(key string, value int) int {
	target, found := os.LookupEnv(key)
	if !found {
		return value
	}

	parsed, err := strconv.Atoi(target)
	if err != nil {
		return value
	}
	return parsed
}

func GetOrPanic(key string) string {
	target, found := os.LookupEnv(key)
	if !found {
		panic(fmt.Sprintf("environment variable %s must be set", key))
	}
	return target
}
