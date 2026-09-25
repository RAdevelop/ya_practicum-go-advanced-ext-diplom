package router

import (
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	uriUserRegister        string = "/api/user/register"
	uriUserLogin           string = "/api/user/login"
	uriUserOrderUpload     string = "/api/user/orders"
	uriUserOrders          string = "/api/user/orders"
	uriUserBalance         string = "/api/user/balance"
	uriUserBalanceWithdraw string = "/api/user/balance/withdraw"
	uriUserWithdrawals     string = "/api/user/withdrawals"
)

func New(handlers *handler.Handlers) http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	routeWithAllowContentTypeApplicationJSON := r.With(middleware.AllowContentType("application/json"))

	routeWithAllowContentTypeApplicationJSON.Post(uriUserRegister, handlers.UserRegister.ServeHTTP)

	routeWithAllowContentTypeApplicationJSON.Post(uriUserLogin, handlers.UserLogin.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("text/plain")).Post(uriUserOrderUpload, handlers.OrderUpload.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Get(uriUserOrders, handlers.Orders.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Get(uriUserBalance, handlers.Balance.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Post(uriUserBalanceWithdraw, handlers.BalanceWithdraw.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Get(uriUserWithdrawals, handlers.BalanceWithdrawals.ServeHTTP)

	return r
}
