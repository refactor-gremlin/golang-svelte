package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// JWTOptions controls how tokens are generated.
type JWTOptions struct {
	Key                      string
	Issuer                   string
	Audience                 string
	AccessTokenLifetimeHours int
}

// Validate ensures all fields are populated and sufficiently strong.
func (o JWTOptions) Validate() error {
	if strings.TrimSpace(o.Key) == "" {
		return errors.New("jwt: key must be provided")
	}
	keyBytes, err := decodeKey(o.Key)
	if err != nil {
		return fmt.Errorf("jwt: invalid key: %w", err)
	}
	if len(keyBytes) < 32 {
		return errors.New("jwt: key must be at least 32 bytes after decoding")
	}

	if strings.TrimSpace(o.Issuer) == "" {
		return errors.New("jwt: issuer must be provided")
	}
	if strings.TrimSpace(o.Audience) == "" {
		return errors.New("jwt: audience must be provided")
	}
	if o.AccessTokenLifetimeHours < 1 || o.AccessTokenLifetimeHours > 168 {
		return errors.New("jwt: access token lifetime must be between 1 and 168 hours")
	}

	return nil
}

// DecodeKey handles both plain text and base64-encoded key formats.
func DecodeKey(key string) ([]byte, error) {
	return decodeKey(key)
}

func decodeKey(key string) ([]byte, error) {
	if strings.HasPrefix(key, "base64:") {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(key, "base64:"))
		if err != nil {
			return nil, err
		}
		return decoded, nil
	}
	return []byte(key), nil
}
