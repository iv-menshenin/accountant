package configstorage

import (
	"os"
	"strconv"
	"time"
)

// EnvInt returns integer environment variable. Returns default if there is no env. Panics on wrong integer format
func EnvInt(envName string, def int) int {
	if val := os.Getenv(envName); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		return i
	}
	return def
}

func EnvString(envName, def string) string {
	if val := os.Getenv(envName); val != "" {
		return val
	}
	return def
}

func EnvDuration(envName string, def time.Duration) time.Duration {
	if val := os.Getenv(envName); val != "" {
		d, err := time.ParseDuration(val)
		if err != nil {
			panic(err)
		}
		return d
	}
	return def
}
