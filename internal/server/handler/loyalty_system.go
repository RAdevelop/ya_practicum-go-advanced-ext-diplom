package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service"
)

// LoyaltySystem - обработка запросов к api программы лояльности
type LoyaltySystem struct {
	appContext     *appcontext.AppContext
	loyaltyManager *service.LoyaltyManager
}

func NewLoyaltySystem(appContext *appcontext.AppContext, loyaltyManager *service.LoyaltyManager) *LoyaltySystem {
	return &LoyaltySystem{
		appContext:     appContext,
		loyaltyManager: loyaltyManager,
	}
}

// UserRegister - Регистрация пользователя
func (ls LoyaltySystem) UserRegister(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement UserRegister", http.StatusNotImplemented)
}

// UserLogin - Аутентификация пользователя
func (ls LoyaltySystem) UserLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement UserLogin", http.StatusNotImplemented)
}

// OrderUpload - Загрузка заказа
func (ls LoyaltySystem) OrderUpload(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement OrderUpload", http.StatusNotImplemented)
}

// Orders - Получение списка загруженных номеров заказов
func (ls LoyaltySystem) Orders(w http.ResponseWriter, r *http.Request) {

	/*
		TODO get from token JWT and convert to uint64?!!
			так же подумать, как сделать получения id пользователя для тестов
	*/
	userID := uint64(1)

	orders, err := ls.loyaltyManager.Orders(userID)
	if err != nil {
		ls.appContext.Logger.Error("Orders", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(orders)

	if err != nil {
		ls.appContext.Logger.Error("Can't write response", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// Balance - Получение текущего баланса пользователя
func (ls LoyaltySystem) Balance(w http.ResponseWriter, r *http.Request) {

	/*
		TODO get from token JWT and convert to uint64?!!
			так же подумать, как сделать получения id пользователя для тестов
	*/
	userID := uint64(1)

	balance, err := ls.loyaltyManager.Balance(userID)
	if err != nil {
		ls.appContext.Logger.Error("Balance", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		ls.appContext.Logger.Error("Can't write response", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// BalanceWithdraw - Запрос на списание средств
func (ls LoyaltySystem) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO implement BalanceWithdraw", http.StatusNotImplemented)
}

// BalanceWithdrawals - Получение информации о выводе средств
func (ls LoyaltySystem) BalanceWithdrawals(w http.ResponseWriter, r *http.Request) {

	/*
		TODO get from token JWT and convert to uint64?!!
			так же подумать, как сделать получения id пользователя для тестов
	*/
	userID := uint64(1)
	withdrawals, err := ls.loyaltyManager.BalanceWithdrawals(userID)

	if err != nil {
		ls.appContext.Logger.Error("BalanceWithdrawals", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		ls.appContext.Logger.Error("Can't write response", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}
