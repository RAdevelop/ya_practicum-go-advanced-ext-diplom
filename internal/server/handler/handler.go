package handler

import (
	"context"
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

type LoyaltyManageable interface {
	OrderUpload(ctx context.Context, userDTO *dto.User, number string) error
	Orders(ctx context.Context, userDTO *dto.User) ([]model.Order, error)
	Balance(ctx context.Context, userDTO *dto.User) (*model.Balance, error)
	BalanceWithdrawals(ctx context.Context, userDTO *dto.User) ([]model.Withdrawal, error)
	BalanceWithdraw(ctx context.Context, userDTO *dto.User, orderNumber string, sum float64) error
	UserRegister(ctx context.Context, userCredentials dto.UserCredentials, jwtSecret string) (string, error)
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
