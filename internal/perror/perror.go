// Package perror - ошибки для проекта
package perror

import "errors"

var (
	ErrInvalidCustomerCredentials      = errors.New("invalid customer credentials")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrOrderAlreadyUploadedByCustomer  = errors.New("order already uploaded by this customer")
	ErrOrderAlreadyUploadedByOther     = errors.New("order already uploaded by another customer")
	ErrBalanceInsufficient             = errors.New("balance insufficient")
	ErrCustomerAlreadyExists           = errors.New("customer already exists")
	ErrCustomerNotFound                = errors.New("customer not found")
	ErrInvalidToken                    = errors.New("invalid token")
	ErrUnexpectedSigningMethodForToken = errors.New("unexpected signing method for token")
	ErrSignToken                       = errors.New("error signing token")
	ErrHashGenerateFromPassword        = errors.New("error generating from password")
)
