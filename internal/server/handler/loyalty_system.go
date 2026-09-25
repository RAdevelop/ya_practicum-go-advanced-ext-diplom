package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/validator"
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
	defer ls.requestBodyClose(r)

	http.Error(w, "TODO implement UserRegister", http.StatusNotImplemented)
}

// UserLogin - Аутентификация пользователя
func (ls LoyaltySystem) UserLogin(w http.ResponseWriter, r *http.Request) {
	defer ls.requestBodyClose(r)

	http.Error(w, "TODO implement UserLogin", http.StatusNotImplemented)
}

// OrderUpload - Загрузка заказа
func (ls LoyaltySystem) OrderUpload(w http.ResponseWriter, r *http.Request) {
	defer ls.requestBodyClose(r)

	if r.Body == nil || r.ContentLength == 0 {
		ls.appContext.Logger.Error("OrderUpload", "body is empty")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		ls.appContext.Logger.Error("OrderUpload", "read body err", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(reqBody))
	if !validator.IsValidLuhn(orderNumber) {
		ls.appContext.Logger.Error("OrderUpload", "invalid orderNumber", orderNumber)
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	/*
		TODO get from token JWT and convert to uint64?!!
			так же подумать, как сделать получения id пользователя для тестов
	*/
	userID := uint64(1)

	err = ls.loyaltyManager.OrderUpload(userID, orderNumber)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	switch {
	case err == nil:
		w.WriteHeader(http.StatusAccepted)

	case errors.Is(err, perror.ErrOrderAlreadyUploadedByUser):
		w.WriteHeader(http.StatusOK)

	case errors.Is(err, perror.ErrOrderAlreadyUploadedByOther):
		ls.appContext.Logger.Error("OrderUpload", "order already uploaded by other user", err)
		http.Error(w, "", http.StatusConflict)

	default:
		ls.appContext.Logger.Error("OrderUpload", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// Orders - Получение списка загруженных номеров заказов
func (ls LoyaltySystem) Orders(w http.ResponseWriter, r *http.Request) {

	defer ls.requestBodyClose(r)

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

	defer ls.requestBodyClose(r)

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
	defer ls.requestBodyClose(r)

	http.Error(w, "TODO implement BalanceWithdraw", http.StatusNotImplemented)
}

// BalanceWithdrawals - Получение информации о выводе средств
func (ls LoyaltySystem) BalanceWithdrawals(w http.ResponseWriter, r *http.Request) {

	defer ls.requestBodyClose(r)

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

func (ls LoyaltySystem) requestBodyClose(r *http.Request) {

	if r.Body == nil {
		return
	}

	err := r.Body.Close()
	if err != nil {
		ls.appContext.Logger.Error("error", "err", err)
	}
}
