package handler

import (
	"context"
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler/middleware"
	statusInner "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/status/inner"
)

//go:generate mockery
type LoyaltyManageable interface {
	// OrderUpload - загрузка заказа
	OrderUpload(ctx context.Context, customerDTO *dto.Customer, number string) error
	// Orders - получение списка загруженных номеров заказов для покупателя
	Orders(ctx context.Context, customerDTO *dto.Customer) ([]dto.Order, error)
	// OrdersAwaitingAccrual - заказы, ожидающие начисления
	OrdersAwaitingAccrual(ctx context.Context, statuses []statusInner.Accrual) ([]dto.Order, error)
	// Balance - получение текущего баланса пользователя
	Balance(ctx context.Context, customerDTO *dto.Customer) (dto.Balance, error)
	// BalanceWithdrawals - получение информации о выводе средств
	BalanceWithdrawals(ctx context.Context, customerDTO *dto.Customer) ([]dto.Withdrawal, error)
	// BalanceWithdraw - списание средств
	BalanceWithdraw(ctx context.Context, customerDTO *dto.Customer, withdraw dto.BalanceWithdraw) error
	// BalanceAccrual - начисление баллов к заказу
	BalanceAccrual(ctx context.Context, accrual dto.Accrual) error
	// UserRegister - регистрация покупателя
	UserRegister(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error)
	// UserLogin - авторизация покупателя
	UserLogin(ctx context.Context, customerCredentials dto.CustomerCredentials, jwtSecret string) (string, error)
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

	/*
		TODO настоятельно рекомендуют добавить сжатие запросов
	*/

	return &Handlers{
		UserRegister:       http.HandlerFunc(ls.UserRegister),
		UserLogin:          http.HandlerFunc(ls.UserLogin),
		OrderUpload:        middleware.Auth(appContext, http.HandlerFunc(ls.OrderUpload)),
		Orders:             middleware.Auth(appContext, http.HandlerFunc(ls.Orders)),
		Balance:            middleware.Auth(appContext, http.HandlerFunc(ls.Balance)),
		BalanceWithdraw:    middleware.Auth(appContext, http.HandlerFunc(ls.BalanceWithdraw)),
		BalanceWithdrawals: middleware.Auth(appContext, http.HandlerFunc(ls.BalanceWithdrawals)),
	}
}
