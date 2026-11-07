package domain

import "errors"

var (
	// ErrUserNotFound is returned when no user matches the lookup.
	ErrUserNotFound = errors.New("user not found")
	// ErrUsernameTaken is returned when creating a duplicate username.
	ErrUsernameTaken = errors.New("username taken")
)
