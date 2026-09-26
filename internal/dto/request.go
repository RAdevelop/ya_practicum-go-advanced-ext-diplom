// Package dto - описание структур для запросов
package dto

type BalanceWithdraw struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}
