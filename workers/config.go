package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// String возвращает строковое значение из env или дефолт
func String(key string, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

// Int возвращает int из env или дефолт
func Int(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if x, err := strconv.Atoi(val); err == nil {
			return x
		}
	}
	return def
}

// StringSlice возвращает []string из env по разделителю или дефолт
func StringSlice(key, def, sep string) []string {
	raw := String(key, def)
	spl := strings.Split(raw, sep)
	for i := range spl {
		spl[i] = strings.TrimSpace(spl[i])
	}
	return spl
}

// Example: config.Must("REDIS_ADDR") для обязательных переменных
func Must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Required env %s not set", key))
	}
	return val
}
