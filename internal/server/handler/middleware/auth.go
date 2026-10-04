package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/jwtoken"
)

const keyUserDTO = "userDTO"
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
		userDTO, err := jwtoken.Parse(token, []byte(appContext.ServerConfig.JWTSecret()))

		if err != nil {
			appContext.Logger.Error("jwtoken.Parse", "error", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(UserDTOPutToCtx(r.Context(), userDTO)))
	})
}

// UserDTOGetFromCtx - получили данные пользователя из контекста
func UserDTOGetFromCtx(ctx context.Context) (*dto.User, bool) {
	user, ok := ctx.Value(keyUserDTO).(*dto.User)
	return user, ok
}

// UserDTOPutToCtx - записали данные пользователя в контекст
func UserDTOPutToCtx(ctx context.Context, user *dto.User) context.Context {
	return context.WithValue(ctx, keyUserDTO, user)
}
