// Package domain holds the auth bounded-context core: users and identity.
package domain

import "time"

// User is the authentication aggregate root.
type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
