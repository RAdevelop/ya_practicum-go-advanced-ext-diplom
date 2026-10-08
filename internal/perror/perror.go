// Package perror - ошибки для проекта
package perror

import "errors"

var (
	ErrCustomerInvalidCredentials = errors.New("invalid customer credentials")
	ErrCustomerInvalidModel       = errors.New("invalid customer model")
	ErrCustomerAlreadyExists      = errors.New("customer already exists")
	ErrCustomerNotFound           = errors.New("customer not found")

	ErrOrderInvalidModel              = errors.New("invalid order")
	ErrOrderNotFound                  = errors.New("order not found")
	ErrOrderAlreadyUploadedByCustomer = errors.New("order already uploaded by this customer")
	ErrOrderAlreadyUploadedByOther    = errors.New("order already uploaded by another customer")
	ErrOrderAccrualAlreadyProcessed   = errors.New("accrual already processed")

	ErrBalanceInsufficient = errors.New("balance insufficient")
	ErrBalanceIncrement    = errors.New("balance increment error")

	ErrTokenInvalid                 = errors.New("invalid token")
	ErrTokenUnexpectedSigningMethod = errors.New("unexpected signing method for token")
	ErrTokenSign                    = errors.New("error signing token")

	ErrHashGenerateFromPassword = errors.New("error generating from password")

	ErrAccrualStatusInvalid = errors.New("invalid accrual status")
	ErrAccrualApply         = errors.New("error accrual apply")
)
