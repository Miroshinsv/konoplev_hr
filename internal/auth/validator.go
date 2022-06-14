package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/meBazil/hr-mvp/internal/config"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidTokenClaims      = fmt.Errorf("invalid token claims")
)

type Validator struct {
	cfg config.JWT
}

func NewValidator(cfg config.JWT) *Validator {
	return &Validator{
		cfg: cfg,
	}
}

func (v *Validator) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, v.keyResolver)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidTokenClaims
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

func (v *Validator) keyResolver(token *jwt.Token) (interface{}, error) {
	// Don't forget to keyResolver the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnexpectedSigningMethod, token.Header["alg"])
	}

	return []byte(v.cfg.Secret), nil
}
