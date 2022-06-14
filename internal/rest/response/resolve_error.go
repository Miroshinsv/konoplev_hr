package response

import (
	"errors"

	"gorm.io/gorm"

	"github.com/meBazil/hr-mvp/internal/auth"
	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

func ResolveError(err error) Error {
	switch {
	case errors.Is(err, auth.ErrDuplicateEmail):
		return NewDuplicateEntryError()
	case errors.Is(err, auth.ErrInvalidCredentials):
		return NewInvalidCredentialsError()
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewNotFoundError()
	case errors.Is(err, auth.ErrInActiveProfile):
		return NewPermissionDeniedError()
	}

	logger.Error("resolved error", err)

	return NewInternalError()
}
