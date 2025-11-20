package etc

import (
	"github.com/dchest/uniuri"
)

// GenerateAPIKey generates a unique API key.
func GenerateAPIKey() (string, error) {
	key := uniuri.NewLenChars(20, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"))
	return key, nil
}

// GenerateAPISecret generates a unique API secret.
func GenerateAPISecret() (string, error) {
	key := uniuri.NewLenChars(30, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"))
	return key, nil
}
