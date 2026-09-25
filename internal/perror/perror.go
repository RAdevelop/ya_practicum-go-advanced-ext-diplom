// Package perror - ошибки для проекта
package perror

import "errors"

var (
	ErrOrderAlreadyUploadedByUser  = errors.New("order already uploaded by this user")
	ErrOrderAlreadyUploadedByOther = errors.New("order already uploaded by another user")
)
