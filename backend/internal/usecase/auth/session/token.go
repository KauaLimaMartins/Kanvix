package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Key(token string) string { return "sess:" + token }

