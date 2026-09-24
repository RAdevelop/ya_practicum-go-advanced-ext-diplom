package handler

import (
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/app_context"
)

type LoyaltySystem struct {
	appContext *app_context.AppContext
}

func NewLoyaltySystem(appContext *app_context.AppContext) *LoyaltySystem {
	return &LoyaltySystem{
		appContext: appContext,
	}
}

func (ls LoyaltySystem) UserRegister(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement UserRegister", http.StatusNotImplemented)
}

func (ls LoyaltySystem) UserLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement UserLogin", http.StatusNotImplemented)
}

func (ls LoyaltySystem) OrderUpload(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement OrderUpload", http.StatusNotImplemented)
}

func (ls LoyaltySystem) Orders(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement Orders", http.StatusNotImplemented)
}

func (ls LoyaltySystem) Balance(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement Balance", http.StatusNotImplemented)
}

func (ls LoyaltySystem) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement BalanceWithdraw", http.StatusNotImplemented)
}

func (ls LoyaltySystem) BalanceWithdrawals(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement BalanceWithdrawals", http.StatusNotImplemented)
}
