package service

import (
	"context"

	"github.com/InternalTransfer/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type AccountRepo interface {
	Create(ctx context.Context, accountID int64, initialBalance decimal.Decimal) error
	GetByID(ctx context.Context, accountID int64) (*model.Account, error)
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, accountID int64) (*model.Account, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, accountID int64, newBalance decimal.Decimal) error
}

type TransactionRepo interface {
	Create(ctx context.Context, tx pgx.Tx, sourceID, destID int64, amount decimal.Decimal) error
}

type TxBeginner interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}
