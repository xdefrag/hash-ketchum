package api

import "errors"

var (
	errAccessDenied = errors.New("Access denied")
	errEmptyLogin   = errors.New("Login required")
	errInvalidLogin = errors.New("Invalid login")
)
