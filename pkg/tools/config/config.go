package config

import (
	"strconv"
)

// WithDefaultString returns value if it not empty string otherwise default.
func WithDefaultString(v, d string) string {
	if v == "" {
		return d
	}

	return v
}

// WithDefaultInt return value if it not 0 otherwise default.
func WithDefaultInt(v, d int) int {
	if v == 0 {
		return d
	}

	return v
}

// Atoi return int value from string or default 0.
func Atoi(v string) int {
	r, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return r
}
