package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler/middleware"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/validator"
)

// jsonBufPool - переиспользуемый буфер для json ответов у хендлеров
var jsonBufPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(nil)
	},
}

// LoyaltySystem - обработка запросов к api программы лояльности
type LoyaltySystem struct {
	appContext     *appcontext.AppContext
	loyaltyManager LoyaltyManageable
}

func NewLoyaltySystem(appContext *appcontext.AppContext, loyaltyManager LoyaltyManageable) *LoyaltySystem {
	return &LoyaltySystem{
		appContext:     appContext,
		loyaltyManager: loyaltyManager,
	}
}

// UserRegister - Регистрация пользователя
func (ls LoyaltySystem) UserRegister(w http.ResponseWriter, r *http.Request) {
	defer ls.requestBodyClose(r)

	var customerCredentials dto.CustomerCredentials
	if err := json.NewDecoder(r.Body).Decode(&customerCredentials); err != nil {
		ls.appContext.Logger.Error("UserRegister", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := validator.IsValidCustomerCredentials(customerCredentials); err != nil {
		ls.appContext.Logger.Warn("UserRegister", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := ls.loyaltyManager.UserRegister(r.Context(), customerCredentials, ls.appContext.ServerConfig.JWTSecret())

	switch {
	case err == nil:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)

	case errors.Is(err, perror.ErrCustomerAlreadyExists):
		http.Error(w, "", http.StatusConflict)

	default:
		ls.appContext.Logger.Error("UserRegister", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// UserLogin - Аутентификация пользователя
func (ls LoyaltySystem) UserLogin(w http.ResponseWriter, r *http.Request) {
	defer ls.requestBodyClose(r)

	var customerCredentials dto.CustomerCredentials
	if err := json.NewDecoder(r.Body).Decode(&customerCredentials); err != nil {
		ls.appContext.Logger.Error("UserLogin", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := validator.IsValidCustomerCredentials(customerCredentials); err != nil {
		ls.appContext.Logger.Warn("UserLogin", "error", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := ls.loyaltyManager.UserLogin(r.Context(), customerCredentials, ls.appContext.ServerConfig.JWTSecret())

	switch {
	case err == nil:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)

	case errors.Is(err, perror.ErrCustomerNotFound):
		http.Error(w, "", http.StatusUnauthorized)

	default:
		ls.appContext.Logger.Error("UserLogin", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// OrderUpload - Загрузка заказа
func (ls LoyaltySystem) OrderUpload(w http.ResponseWriter, r *http.Request) {
	defer ls.requestBodyClose(r)

	reqBody, hasError := ls.requestBodyGet(w, r, "OrderUpload")
	if hasError {
		return
	}

	orderNumber := strings.TrimSpace(string(reqBody))

	if !validator.IsValidLuhn(orderNumber) {
		ls.appContext.Logger.Error("OrderUpload", "invalid orderNumber", orderNumber)
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	customerDTO, ok := middleware.CustomerGetFromCtx(r.Context())
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	err := ls.loyaltyManager.OrderUpload(r.Context(), customerDTO, orderNumber)

	responseSetHeaderContentTypeTextPlain(w)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusAccepted)

	case errors.Is(err, perror.ErrOrderAlreadyUploadedByCustomer):
		w.WriteHeader(http.StatusOK)

	case errors.Is(err, perror.ErrOrderAlreadyUploadedByOther):
		ls.appContext.Logger.Error("OrderUpload", "order already uploaded by other customer", err)
		http.Error(w, "", http.StatusConflict)

	default:
		ls.appContext.Logger.Error("OrderUpload", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// Orders - Получение списка загруженных номеров заказов
func (ls LoyaltySystem) Orders(w http.ResponseWriter, r *http.Request) {

	defer ls.requestBodyClose(r)

	customerDTO, ok := middleware.CustomerGetFromCtx(r.Context())

	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	orders, err := ls.loyaltyManager.Orders(r.Context(), customerDTO)
	if err != nil {
		ls.appContext.Logger.Error("Orders", "err", err)
		if errors.Is(err, perror.ErrOrderNotFound) {
			err = nil
		}
	}

	ls.responseJSON(w, orders, len(orders) == 0, err)
}

// Balance - Получение текущего баланса пользователя
func (ls LoyaltySystem) Balance(w http.ResponseWriter, r *http.Request) {

	defer ls.requestBodyClose(r)

	customerDTO, ok := middleware.CustomerGetFromCtx(r.Context())

	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	balance, err := ls.loyaltyManager.Balance(r.Context(), customerDTO)
	if err != nil {
		ls.appContext.Logger.Error("Balance", "err", err)
		if errors.Is(err, perror.ErrBalanceCustomerNotFound) {
			err = nil
		}
	}

	ls.responseJSON(w, balance, false, err)
}

// BalanceWithdraw - Запрос на списание средств
func (ls LoyaltySystem) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	defer ls.requestBodyClose(r)

	reqBody, hasError := ls.requestBodyGet(w, r, "BalanceWithdraw")
	if hasError {
		return
	}

	var balanceWithdraw dto.BalanceWithdraw
	err := json.Unmarshal(reqBody, &balanceWithdraw)
	if err != nil {
		ls.appContext.Logger.Error("Can't parse balanceWithdraw", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if !validator.IsValidLuhn(balanceWithdraw.OrderNumber) {
		ls.appContext.Logger.Error("BalanceWithdraw", "invalid orderNumber", balanceWithdraw)
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	if balanceWithdraw.Sum <= 0 {
		ls.appContext.Logger.Error("BalanceWithdraw", "invalid sum", balanceWithdraw)
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	customerDTO, ok := middleware.CustomerGetFromCtx(r.Context())

	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	err = ls.loyaltyManager.BalanceWithdraw(r.Context(), customerDTO, balanceWithdraw)
	responseSetHeaderContentTypeTextPlain(w)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, perror.ErrBalanceInsufficient):
		ls.appContext.Logger.Error("BalanceWithdraw", "ErrBalanceInsufficient", err)
		http.Error(w, "", http.StatusPaymentRequired)
	case errors.Is(err, perror.ErrOrderNotFound):
		ls.appContext.Logger.Error("BalanceWithdraw", "ErrOrderNotFound", err)
		http.Error(w, "", http.StatusUnprocessableEntity)
	default:
		ls.appContext.Logger.Error("BalanceWithdraw", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// BalanceWithdrawals - Получение информации о выводе средств
func (ls LoyaltySystem) BalanceWithdrawals(w http.ResponseWriter, r *http.Request) {

	defer ls.requestBodyClose(r)

	customerDTO, ok := middleware.CustomerGetFromCtx(r.Context())

	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	withdrawals, err := ls.loyaltyManager.BalanceWithdrawals(r.Context(), customerDTO)

	if err != nil {
		ls.appContext.Logger.Error("BalanceWithdrawals", "err", err)
		if errors.Is(err, perror.ErrWithdrawalNotFound) {
			err = nil
		}
	}

	ls.responseJSON(w, withdrawals, len(withdrawals) == 0, err)
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

func (ls LoyaltySystem) requestBodyGet(w http.ResponseWriter, r *http.Request, methodName string) ([]byte, bool) {
	if r.Body == nil || r.ContentLength == 0 {
		ls.appContext.Logger.Error(methodName, "body is empty")
		http.Error(w, "", http.StatusBadRequest)
		return nil, true
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ls.appContext.Logger.Error(methodName, "read body err", err)
		http.Error(w, "", http.StatusBadRequest)
		return nil, true
	}

	return body, false
}

func responseSetHeaderContentTypeTextPlain(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (ls LoyaltySystem) responseJSON(w http.ResponseWriter, respData any, isRespDataEmpty bool, err error) {

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if isRespDataEmpty {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBuf := jsonBufPool.Get().(*bytes.Buffer)
	jsonBuf.Reset()
	defer jsonBufPool.Put(jsonBuf)

	if err = json.NewEncoder(jsonBuf).Encode(respData); err != nil {
		ls.appContext.Logger.Error("Can't write response", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBuf.Bytes())
}
