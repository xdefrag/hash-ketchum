package api

import "context"

// Credentials contains login for logining on server.
type Credentials struct {
	Login string
}

// GetRequestMetadata extracts metadata from req.
func (c Credentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login": c.Login,
	}, nil
}

// RequireTransportSecurity setting.
func (c Credentials) RequireTransportSecurity() bool {
	return false
}
