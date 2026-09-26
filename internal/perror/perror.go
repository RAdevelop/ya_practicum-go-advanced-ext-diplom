// Package perror - ошибки для проекта
package perror

import "errors"

var (
	ErrOrderNotFound               = errors.New("order not found")
	ErrOrderAlreadyUploadedByUser  = errors.New("order already uploaded by this user")
	ErrOrderAlreadyUploadedByOther = errors.New("order already uploaded by another user")
	ErrBalanceInsufficient         = errors.New("balance insufficient")
)
