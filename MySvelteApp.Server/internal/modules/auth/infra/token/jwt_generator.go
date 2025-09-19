package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authapp "mysvelteapp/server_new/internal/modules/auth/app"
	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

var _ authapp.TokenGenerator = (*JWTTokenGenerator)(nil)

// JWTTokenGenerator implements TokenGenerator using github.com/golang-jwt/jwt/v5.
type JWTTokenGenerator struct {
	options    JWTOptions
	signingKey []byte
}

// NewJWTTokenGenerator validates the provided options and prepares a generator instance.
func NewJWTTokenGenerator(options JWTOptions) (*JWTTokenGenerator, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	keyBytes, err := DecodeKey(options.Key)
	if err != nil {
		return nil, fmt.Errorf("decode key: %w", err)
	}

	return &JWTTokenGenerator{
		options:    options,
		signingKey: keyBytes,
	}, nil
}

// GenerateToken produces a signed JWT for the supplied user entity.
func (g *JWTTokenGenerator) GenerateToken(user *authdomain.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user must not be nil")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(g.options.AccessTokenLifetimeHours) * time.Hour)

	claims := authClaims{
		Username: user.Username,
		NameID:   fmt.Sprintf("%d", user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			Issuer:    g.options.Issuer,
			Audience:  []string{g.options.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(g.signingKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

type authClaims struct {
	Username string `json:"name"`
	NameID   string `json:"nameid"`
	jwt.RegisteredClaims
}
