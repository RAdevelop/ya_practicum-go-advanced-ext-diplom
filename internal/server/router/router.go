package router

import (
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	uriApiUserRegister        string = "/api/user/register"
	uriApiUserLogin           string = "/api/user/login"
	uriApiUserOrderUpload     string = "/api/user/orders"
	uriApiUserOrders          string = "/api/user/orders"
	uriApiUserBalance         string = "/api/user/balance"
	uriApiUserBalanceWithdraw string = "/api/user/balance/withdraw"
	uriApiUserWithdrawals     string = "/api/user/withdrawals"
)

func New(handlers *handler.Handlers) http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	routeWithAllowContentTypeApplicationJSON := r.With(middleware.AllowContentType("application/json"))

	routeWithAllowContentTypeApplicationJSON.Post(uriApiUserRegister, handlers.UserRegister.ServeHTTP)

	routeWithAllowContentTypeApplicationJSON.Post(uriApiUserLogin, handlers.UserLogin.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("text/plain")).Post(uriApiUserOrderUpload, handlers.OrderUpload.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Get(uriApiUserOrders, handlers.Orders.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Get(uriApiUserBalance, handlers.Balance.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Post(uriApiUserBalanceWithdraw, handlers.BalanceWithdrawals.ServeHTTP)

	//TODO implement + проверка аутентификации
	routeWithAllowContentTypeApplicationJSON.Get(uriApiUserWithdrawals, handlers.BalanceWithdrawals.ServeHTTP)

	return r
}
