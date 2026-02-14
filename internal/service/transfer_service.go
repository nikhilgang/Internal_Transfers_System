package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"github.com/InternalTransfer/internal/apperror"
)

type TransferService struct {
	accountRepo       AccountRepo
	transactionRepo   TransactionRepo
	txBeginner        TxBeginner
	logger            *slog.Logger
	maxTransferAmount decimal.Decimal
}

var minTransferAmount = decimal.NewFromInt(1)

const maxRetries = 3

func NewTransferService(
	accountRepo AccountRepo,
	transactionRepo TransactionRepo,
	txBeginner TxBeginner,
	logger *slog.Logger,
	maxTransferAmount int64,
) *TransferService {
	return &TransferService{
		accountRepo:       accountRepo,
		transactionRepo:   transactionRepo,
		txBeginner:        txBeginner,
		logger:            logger,
		maxTransferAmount: decimal.NewFromInt(maxTransferAmount),
	}
}

func (s *TransferService) Transfer(ctx context.Context, sourceID, destID int64, amount decimal.Decimal) error {
	if sourceID <= 0 || destID <= 0 {
		return &apperror.ErrValidation{Message: "Please provide valid account numbers"}
	}
	if sourceID == destID {
		return &apperror.ErrValidation{Message: "Cannot transfer to the same account. Please choose a different destination account"}
	}
	if amount.LessThan(minTransferAmount) {
		return &apperror.ErrValidation{Message: fmt.Sprintf("Transfer amount must be at least $%s", minTransferAmount)}
	}
	if amount.GreaterThan(s.maxTransferAmount) {
		return &apperror.ErrValidation{Message: fmt.Sprintf("Transfer amount cannot exceed $%s", s.maxTransferAmount)}
	}
	if amount.Exponent() < -2 {
		return &apperror.ErrValidation{Message: "Transfer amount can only have up to 2 decimal places (e.g., 10.50)"}
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		lastErr = s.executeTransfer(ctx, sourceID, destID, amount)
		if lastErr == nil {
			return nil
		}

		var pgErr *pgconn.PgError
		if errors.As(lastErr, &pgErr) && pgErr.Code == "40001" {
			s.logger.Warn("serialization failure, retrying", "attempt", i+1)
			continue
		}

		return lastErr
	}

	return fmt.Errorf("unable to complete transfer due to high system load. Please try again in a few moments")
}

func (s *TransferService) executeTransfer(ctx context.Context, sourceID, destID int64, amount decimal.Decimal) error {
	tx, err := s.txBeginner.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// lock accounts in consistent order to avoid deadlocks
	firstID, secondID := sourceID, destID
	if sourceID > destID {
		firstID, secondID = destID, sourceID
	}

	first, err := s.accountRepo.GetByIDForUpdate(ctx, tx, firstID)
	if err != nil {
		return err
	}
	second, err := s.accountRepo.GetByIDForUpdate(ctx, tx, secondID)
	if err != nil {
		return err
	}

	sourceAccount, destAccount := first, second
	if sourceID != firstID {
		sourceAccount, destAccount = second, first
	}

	if sourceAccount.Balance.LessThan(amount) {
		return &apperror.ErrInsufficientBalance{AccountID: sourceID}
	}

	newSourceBal := sourceAccount.Balance.Sub(amount)
	newDestBal := destAccount.Balance.Add(amount)

	if err = s.accountRepo.UpdateBalance(ctx, tx, sourceID, newSourceBal); err != nil {
		return err
	}
	if err = s.accountRepo.UpdateBalance(ctx, tx, destID, newDestBal); err != nil {
		return err
	}

	if err = s.transactionRepo.Create(ctx, tx, sourceID, destID, amount); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transfer: %w", err)
	}

	s.logger.Info("transfer completed", "source", sourceID, "destination", destID, "amount", amount.String())
	return nil
}
