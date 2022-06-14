package auth

import (
	"github.com/dgrijalva/jwt-go"

	helper "github.com/meBazil/hr-mvp/internal/rest/helpers"
)

const (
	RefreshToken Type = iota + 1
	AccessToken

	UserSubject = "user"
)

type Type uint

type Claims struct {
	jwt.StandardClaims

	UserID  uint `json:"uid"`
	Version uint `json:"ver"`
	Type    Type `json:"typ"`
}

func (c *Claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.Type == 0 {
		return helper.ErrInvalidTokenType
	}

	return nil
}
