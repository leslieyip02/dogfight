package env

import (
	"fmt"
	"os"
)

func GetOrDefault(key string, value string) string {
	target, found := os.LookupEnv(key)
	if !found {
		return value
	}
	return target
}

func GetOrPanic(key string) string {
	target, found := os.LookupEnv(key)
	if !found {
		panic(fmt.Sprintf("environment variable %s must be set", key))
	}
	return target
}
