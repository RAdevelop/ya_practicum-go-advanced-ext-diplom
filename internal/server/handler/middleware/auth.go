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

func Auth(appContext appcontext.AppContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			http.Error(w, "", http.StatusUnauthorized)
		}
		token := strings.TrimPrefix(header, prefix)

		userDTO, err := jwtoken.Parse(token, []byte(appContext.ServerConfig.JWTSecret()))

		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r.WithContext(UserDTOPutToCtx(r.Context(), userDTO)))
	})
}

func UserDTOGetFromCtx(ctx context.Context) (*dto.User, bool) {
	user, ok := ctx.Value(keyUserDTO).(*dto.User)
	return user, ok
}
func UserDTOPutToCtx(ctx context.Context, user *dto.User) context.Context {
	return context.WithValue(ctx, keyUserDTO, user)
}
