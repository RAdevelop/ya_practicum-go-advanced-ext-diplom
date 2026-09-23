package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	r.Post("/api/user/register", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement
		http.Error(w, "TODO implement", http.StatusNotFound)
	})
	r.Post("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement
		http.Error(w, "TODO implement", http.StatusNotFound)
	})

	r.Post("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement + проверка аутентификации
		http.Error(w, "TODO implement", http.StatusNotFound)
	})

	r.Get("/api/user/orders", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement + проверка аутентификации
		http.Error(w, "TODO implement", http.StatusNotFound)
	})

	r.Get("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement + проверка аутентификации
		http.Error(w, "TODO implement", http.StatusNotFound)
	})

	r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement + проверка аутентификации
		http.Error(w, "TODO implement", http.StatusNotFound)
	})

	r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {
		//TODO implement + проверка аутентификации
		http.Error(w, "TODO implement", http.StatusNotFound)
	})

	return r
}
