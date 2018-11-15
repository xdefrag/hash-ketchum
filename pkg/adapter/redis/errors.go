package redis

import "fmt"
import "strings"

// ErrEmptyConfigValues error with listing of empty keys.
type ErrEmptyConfigValues struct {
	keys []string
}

// NewEmptyConfigValues build new ErrEmptyConfigValues with empty keys.
func NewEmptyConfigValues(keys []string) ErrEmptyConfigValues {
	return ErrEmptyConfigValues{keys}
}

// Error message.
func (e ErrEmptyConfigValues) Error() string {
	return fmt.Sprintf("No config value in key: %s", strings.Join(e.keys, ", "))
}
