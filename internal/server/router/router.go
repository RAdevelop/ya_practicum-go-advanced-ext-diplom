package router

import (
	"net/http"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(handlers *handler.Handlers) http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	r.With(middleware.AllowContentType("application/json")).
		Post("/api/user/register", handlers.UserRegister.ServeHTTP)

	r.With(middleware.AllowContentType("application/json")).
		Post("/api/user/login", handlers.UserLogin.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("text/plain")).
		Post("/api/user/orders", handlers.OrderUpload.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("application/json")).
		Get("/api/user/orders", handlers.Orders.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("application/json")).
		Get("/api/user/balance", handlers.Balance.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("application/json")).
		Post("/api/user/balance/withdraw", handlers.BalanceWithdrawals.ServeHTTP)

	//TODO implement + проверка аутентификации
	r.With(middleware.AllowContentType("application/json")).
		Get("/api/user/withdrawals", handlers.BalanceWithdrawals.ServeHTTP)

	return r
}
