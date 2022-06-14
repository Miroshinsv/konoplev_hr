package helper

import "errors"

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidTokenType    = errors.New("invalid token type")
	ErrInvalidTokenVersion = errors.New("invalid token version")
)
