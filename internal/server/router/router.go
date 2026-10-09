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

	r.With(middleware.AllowContentType("text/plain")).Post(uriUserOrderUpload, handlers.OrderUpload.ServeHTTP)

	routeWithAllowContentTypeApplicationJSON.Get(uriUserOrders, handlers.Orders.ServeHTTP)

	routeWithAllowContentTypeApplicationJSON.Get(uriUserBalance, handlers.Balance.ServeHTTP)

	routeWithAllowContentTypeApplicationJSON.Post(uriUserBalanceWithdraw, handlers.BalanceWithdraw.ServeHTTP)

	routeWithAllowContentTypeApplicationJSON.Get(uriUserWithdrawals, handlers.BalanceWithdrawals.ServeHTTP)

	return r
}
