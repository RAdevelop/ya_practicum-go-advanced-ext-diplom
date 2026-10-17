package jwtoken

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateParse_RoundTrip(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")
	customer := dto.Customer{ID: 42, Login: "alice"}

	token, err := Generate(customer, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsed, err := Parse(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, customer.ID, parsed.ID)
	assert.Equal(t, customer.Login, parsed.Login)
}

func TestParse_WrongSecret(t *testing.T) {
	secret1 := []byte("super-secret-key-at-least-32-bytes-long")
	secret2 := []byte("another-secret-key-at-least-32-bytes-long")

	token, err := Generate(dto.Customer{ID: 1, Login: "alice"}, secret1)
	assert.NoError(t, err)

	_, err = Parse(token, secret2)
	assert.ErrorIs(t, err, perror.ErrTokenInvalid)
}

func TestParse_MalformedToken(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")

	cases := []string{
		"",
		"not-a-jwt",
		"a.b.c",
		"eyJhbGciOiJIUzI1NiJ9.invalid.signature",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			_, err := Parse(tc, secret)
			assert.ErrorIs(t, err, perror.ErrTokenInvalid)
		})
	}
}

func TestParse_ExpiredToken(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")

	// создаём вручную токен с exp в прошлом
	claims := dto.CustomerJWT{
		CustomerLogin: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	assert.NoError(t, err)

	_, err = Parse(signed, secret)
	assert.ErrorIs(t, err, perror.ErrTokenInvalid)
}

func TestParse_NoneAlgorithm(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")

	// вручную клеим токен с alg=none
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"1","customer_login":"alice"}`))
	token := header + "." + payload + "."

	_, err := Parse(token, secret)
	assert.Error(t, err)
}

func TestParse_HS512Rejected(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")

	claims := dto.CustomerJWT{
		CustomerLogin: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // ← HS512
	signed, err := token.SignedString(secret)
	assert.NoError(t, err)

	_, err = Parse(signed, secret)
	assert.Error(t, err) // ожидаем только HS256
}

func TestGenerateParse_EmptyLogin(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")
	customer := dto.Customer{ID: 1, Login: ""}

	token, err := Generate(customer, secret)
	assert.NoError(t, err)

	parsed, err := Parse(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), parsed.ID)
	assert.Equal(t, "", parsed.Login)
}

func TestParse_InvalidSubject(t *testing.T) {
	secret := []byte("super-secret-key-at-least-32-bytes-long")

	claims := dto.CustomerJWT{
		CustomerLogin: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "not-a-number",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	assert.NoError(t, err)

	_, err = Parse(signed, secret)
	assert.ErrorIs(t, err, perror.ErrTokenInvalid)
}
