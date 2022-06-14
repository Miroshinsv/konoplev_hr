package auth

import (
	base "github.com/dgrijalva/jwt-go"

	"github.com/meBazil/hr-mvp/internal/config"
)

const IssuerName = "taxes"

type Issuer struct {
	cfg config.JWT
}

func NewIssuer(cfg config.JWT) *Issuer {
	return &Issuer{
		cfg: cfg,
	}
}

func (i *Issuer) Issue(claims *Claims) (string, error) {
	token := base.NewWithClaims(base.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(i.cfg.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
