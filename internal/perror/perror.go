// Package perror - ошибки для проекта
package perror

import "errors"

var (
	ErrInvalidUserCredentials          = errors.New("invalid user credentials")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrOrderAlreadyUploadedByUser      = errors.New("order already uploaded by this user")
	ErrOrderAlreadyUploadedByOther     = errors.New("order already uploaded by another user")
	ErrBalanceInsufficient             = errors.New("balance insufficient")
	ErrUserAlreadyExists               = errors.New("login already exists")
	ErrUserNotFound                    = errors.New("user not found")
	ErrInvalidToken                    = errors.New("invalid token")
	ErrUnexpectedSigningMethodForToken = errors.New("unexpected signing method for token")
	ErrSignToken                       = errors.New("error signing token")
	ErrHashGenerateFromPassword        = errors.New("error generating from password")
)
