package handler

import (
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service"
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

func New(appContext *appcontext.AppContext, loyaltyManager *service.LoyaltyManager) *Handlers {

	ls := NewLoyaltySystem(appContext, loyaltyManager)

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
