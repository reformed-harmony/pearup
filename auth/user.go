package auth

import (
	"io"
)

// User stores user information provided by a provider.
type User struct {
	ID      string
	Name    string
	Email   string
	Picture io.ReadCloser
}
