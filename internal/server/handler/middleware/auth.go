package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
)

type contextKey struct{}

var keyCustomerDTO = contextKey{}

const prefix = "Bearer "

func Auth(appContext *appcontext.AppContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(header, prefix) {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, prefix)
		customer, err := jwtoken.Parse(token, []byte(appContext.ServerConfig.JWTSecret()))

		if err != nil {
			appContext.Logger.Error("jwtoken.Parse", "error", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(CustomerPutToCtx(r.Context(), customer)))
	})
}

// CustomerGetFromCtx - получили данные пользователя из контекста
func CustomerGetFromCtx(ctx context.Context) (*dto.Customer, bool) {
	customer, ok := ctx.Value(keyCustomerDTO).(*dto.Customer)
	return customer, ok
}

// CustomerPutToCtx - записали данные пользователя в контекст
func CustomerPutToCtx(ctx context.Context, customer *dto.Customer) context.Context {
	return context.WithValue(ctx, keyCustomerDTO, customer)
}
