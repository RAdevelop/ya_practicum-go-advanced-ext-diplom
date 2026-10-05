package jwtoken

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// TokenTTL — время жизни токена.
	tokenTTL = 24 * time.Hour

	// SigningMethod — алгоритм подписи.
	// Используется и при создании, и при проверке.
	signingAlg = "HS256"
)

/*
Generate - создаёт подписанный JWT для пользователя.

- CustomerID попадает в Subject, login — в кастомный клейм.
- Секрет должен быть криптостойким (>= 32 байта для HS256).
*/
func Generate(customer dto.Customer, secret []byte) (string, error) {
	now := time.Now()

	claims := dto.CustomerJWT{
		CustomerLogin: customer.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			// Идентификатор субъекта, к которому относится токен. Обычно — ID пользователя
			Subject: strconv.FormatUint(customer.ID, 10),
			// Когда выдан
			IssuedAt: jwt.NewNumericDate(now),
			// Когда истекает
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("%w: %w", perror.ErrSignToken, err)
	}

	return signed, nil
}

/*
Parse проверяет подпись и срок действия токена, возвращает объект пользователя.

Проверяет:
  - подпись HMAC (защита от algorithm confusion),
  - exp (автоматически библиотекой),
  - nbf (автоматически библиотекой, если задан).
*/
func Parse(tokenString string, secret []byte) (*dto.Customer, error) {
	claims := &dto.CustomerJWT{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			// Проверяем, что алгоритм — HMAC.
			// Защита от подмены alg на "none" или асимметричный.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", perror.ErrUnexpectedSigningMethodForToken, t.Header["alg"])
			}
			return secret, nil
		},
		// Явно разрешаем только HS256.
		jwt.WithValidMethods([]string{signingAlg}),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", perror.ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, perror.ErrInvalidToken
	}

	customerID, err := claims.CustomerID()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", perror.ErrInvalidToken, err)
	}

	return &dto.Customer{
		ID:    customerID,
		Login: claims.CustomerLogin,
	}, nil
}
