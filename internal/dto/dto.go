package dto

import "github.com/shopspring/decimal"

type CreateAccountRequest struct {
	AccountID      int64           `json:"account_id"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}

type AccountResponse struct {
	AccountID int64           `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}

type CreateTransactionRequest struct {
	SourceAccountID      int64           `json:"source_account_id"`
	DestinationAccountID int64           `json:"destination_account_id"`
	Amount               decimal.Decimal `json:"amount"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
