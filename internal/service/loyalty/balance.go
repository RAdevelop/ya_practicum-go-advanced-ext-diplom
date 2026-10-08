package loyalty

import (
	"context"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/dto"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/model"
)

// Balance - получение текущего баланса пользователя
func (lm *Manager) Balance(ctx context.Context, customerDTO *dto.Customer) (dto.Balance, error) {
	balanceModel, err := lm.storage.BalanceByCustomerID(ctx, customerDTO.ID)
	if err != nil {
		return dto.Balance{}, err
	}

	balance := dto.Balance{
		Current:   balanceModel.Current,
		Withdrawn: balanceModel.Withdrawn,
	}

	balanceModel = nil
	return balance, nil
}

// BalanceWithdrawals - получение информации о выводе средств
func (lm *Manager) BalanceWithdrawals(ctx context.Context, customerDTO *dto.Customer) ([]dto.Withdrawal, error) {

	withdrawalsModel, err := lm.storage.BalanceWithdrawalsByCustomerID(ctx, customerDTO.ID)
	if err != nil {
		return nil, err
	}

	withdrawals := make([]dto.Withdrawal, 0, len(withdrawalsModel))
	for _, withdrawal := range withdrawalsModel {
		withdrawals = append(withdrawals, dto.Withdrawal{
			Order:       withdrawal.Order,
			Sum:         withdrawal.Sum,
			ProcessedAt: withdrawal.ProcessedAt,
		})
	}

	withdrawalsModel = nil
	return withdrawals, nil
}

// BalanceWithdraw - списание средств
func (lm *Manager) BalanceWithdraw(ctx context.Context, customerDTO *dto.Customer, withdraw dto.BalanceWithdraw) error {

	withdrawal := model.Withdrawal{
		CustomerID: customerDTO.ID,
		Order:      withdraw.OrderNumber,
		Sum:        withdraw.Sum,
	}

	return lm.storage.BalanceWithdraw(ctx, withdrawal)
}

// BalanceAccrual - начисление баллов к заказу
func (lm *Manager) BalanceAccrual(ctx context.Context, accrual dto.Accrual) error {

	accrualModel := model.Accrual{
		Order:   accrual.Order,
		Status:  accrual.Status,
		Accrual: accrual.Accrual,
	}

	return lm.storage.BalanceAccrual(ctx, accrualModel)
}
