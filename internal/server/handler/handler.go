package handler

import (
	"context"
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

type LoyaltyManageable interface {
	OrderUpload(ctx context.Context, userID uint64, number string) error
	Orders(ctx context.Context, userID uint64) ([]model.Order, error)
	Balance(ctx context.Context, userID uint64) (*model.Balance, error)
	BalanceWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error)
	BalanceWithdraw(ctx context.Context, userID uint64, orderNumber string, sum float64) error
}

type Handlers struct {
	UserRegister       http.Handler
	UserLogin          http.Handler
	OrderUpload        http.Handler
	Orders             http.Handler
	Balance            http.Handler
	BalanceWithdraw    http.Handler
	BalanceWithdrawals http.Handler
}

func New(appContext *appcontext.AppContext, loyaltyManager LoyaltyManageable) *Handlers {

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
