package database

import (
	"context"
	"testing"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/perror"
	"github.com/stretchr/testify/assert"
)

func assertCustomer(t *testing.T, expectedCustomer *model.Customer, actualCustomer *model.Customer) {
	t.Helper()

	assert.Equalf(t, expectedCustomer.Login, actualCustomer.Login, "expectedCustomer: %+v", expectedCustomer)
	assert.Equalf(t, expectedCustomer.PasswordHash, actualCustomer.PasswordHash, "expectedCustomer: %+v", expectedCustomer)
	assert.Greaterf(t, actualCustomer.ID, uint64(0), "expectedCustomer: %+v", expectedCustomer)
	assert.WithinDurationf(t, time.Now(), actualCustomer.CreatedAt, time.Minute, "expectedCustomer: %+v", expectedCustomer)
}

func TestStorage_CustomerCreateErrCustomerAlreadyExists(t *testing.T) {

	storage := setUpStorage(t)

	customerAlreadyExists := &model.Customer{
		Login:        "john",
		PasswordHash: "PasswordHash",
	}
	customerNew := &model.Customer{
		Login:        "john",
		PasswordHash: "PasswordHash2",
	}

	t.Run("ErrCustomerAlreadyExists", func(t *testing.T) {

		err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {
			_, err := storage.CustomerCreate(ctx, customerAlreadyExists)
			assert.NoError(t, err)

			customer, err := storage.CustomerCreate(ctx, customerNew)
			assert.ErrorIs(t, err, perror.ErrCustomerAlreadyExists)
			assert.Nilf(t, customer, "created customer: %+v", customer)

			return errForTransactionRollback
		})

		assert.Error(t, err)
	})
}

func TestStorage_CustomerCreate(t *testing.T) {

	type given struct {
		customer *model.Customer
	}

	type want struct {
		err error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				customer: &model.Customer{
					Login:        "john",
					PasswordHash: "PasswordHash",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "error ErrInvalidCustomerModel is empty",
			given: given{
				customer: &model.Customer{},
			},
			want: want{
				err: perror.ErrCustomerInvalidModel,
			},
		},
		{
			name: "error ErrInvalidCustomerModel is nil",
			given: given{
				customer: nil,
			},
			want: want{
				err: perror.ErrCustomerInvalidModel,
			},
		},
	}

	storage := setUpStorage(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {
				customer, err := storage.CustomerCreate(ctx, tt.given.customer)

				assert.Equalf(t, tt.want.err, err, "given: %+v", tt.given)

				if err == nil {
					assertCustomer(t, tt.given.customer, customer)
				} else {
					assert.Nilf(t, customer, "given: %+v", tt.given)
				}

				return errForTransactionRollback
			})

			assert.Error(t, err)
		})
	}
}

func TestStorage_CustomerFindByLogin(t *testing.T) {

	type given struct {
		customer *model.Customer
		login    string
	}
	type want struct {
		err error
	}

	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "success",
			given: given{
				customer: &model.Customer{
					Login:        "john",
					PasswordHash: "PasswordHash",
				},
				login: "john",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "ErrCustomerNotFound",
			given: given{
				customer: &model.Customer{
					Login:        "john",
					PasswordHash: "PasswordHash",
				},
				login: "not-found",
			},
			want: want{
				err: perror.ErrCustomerNotFound,
			},
		},
		{
			name: "ErrCustomerNotFound with empty login",
			given: given{
				customer: &model.Customer{
					Login:        "john",
					PasswordHash: "PasswordHash",
				},
				login: "",
			},
			want: want{
				err: perror.ErrCustomerNotFound,
			},
		},
	}

	storage := setUpStorage(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := storage.DB.RunInTransaction(t.Context(), func(ctx context.Context) error {

				customerCreated, err := storage.CustomerCreate(ctx, tt.given.customer)
				assert.NoError(t, err)

				customerFound, err := storage.CustomerFindByLogin(ctx, tt.given.login)
				assert.ErrorIs(t, err, tt.want.err)

				if tt.want.err == nil {
					assertCustomer(t, customerCreated, customerFound)
				} else {
					assert.Nilf(t, customerFound, "customerFound: %+v", customerFound)
				}

				return errForTransactionRollback
			})
			assert.Error(t, err)
		})
	}
}
