package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/InternalTransfer/internal/apperror"
	"github.com/InternalTransfer/internal/model"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	accountRepo AccountRepo
	logger      *slog.Logger
}

func NewAccountService(accountRepo AccountRepo, logger *slog.Logger) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

func (s *AccountService) Create(ctx context.Context, accountID int64, initialBalance decimal.Decimal) error {
	if accountID <= 0 {
		return &apperror.ErrValidation{Message: "Please provide a valid account number"}
	}
	if initialBalance.IsNegative() {
		return &apperror.ErrValidation{Message: "Initial balance cannot be negative"}
	}
	if initialBalance.Exponent() < -2 {
		return &apperror.ErrValidation{Message: "Initial balance can only have up to 2 decimal places (e.g., 100.50)"}
	}

	if err := s.accountRepo.Create(ctx, accountID, initialBalance); err != nil {
		return fmt.Errorf("creating account: %w", err)
	}

	s.logger.Info("account created", "account_id", accountID)
	return nil
}

func (s *AccountService) GetByID(ctx context.Context, accountID int64) (*model.Account, error) {
	if accountID <= 0 {
		return nil, &apperror.ErrValidation{Message: "Please provide a valid account number"}
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("fetching account: %w", err)
	}

	return account, nil
}
