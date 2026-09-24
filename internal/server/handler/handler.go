package handler

import (
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/app_context"
)

type Handlers struct {
	UserRegister       http.Handler
	UserLogin          http.Handler
	OrderUpload        http.Handler
	Orders             http.Handler
	Balance            http.Handler
	BalanceWithdraw    http.Handler
	BalanceWithdrawals http.Handler
}

func New(appContext *app_context.AppContext) *Handlers {

	ls := NewLoyaltySystem(appContext)

	return &Handlers{
		UserRegister:       http.HandlerFunc(ls.UserRegister),
		UserLogin:          http.HandlerFunc(ls.UserLogin),
		OrderUpload:        http.HandlerFunc(ls.OrderUpload),
		Orders:             http.HandlerFunc(ls.Orders),
		Balance:            http.HandlerFunc(ls.Balance),
		BalanceWithdraw:    http.HandlerFunc(ls.BalanceWithdraw),
		BalanceWithdrawals: http.HandlerFunc(ls.BalanceWithdrawals),
	}
}
