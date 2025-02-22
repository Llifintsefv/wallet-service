package model

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Uuid          string
	OperationType string
	Amount        decimal.Decimal
}

type TransactionRequest struct {
	Uuid          string `json:"uuid"`
	OperationType string `json:"operationType"`
	Amount        string `json:"amount"`
}

const (
	TransactionDeposit  = "DEPOSIT"
	TransactionWithdraw = "WITHDRAW"
)

func ValidateTransaction(req Transaction) error {
	if req.Uuid == "" {
		return fmt.Errorf("uuid is required")
	}
	if req.OperationType != TransactionDeposit && req.OperationType != TransactionWithdraw {
		return fmt.Errorf("invalid operation type: %s", req.OperationType)
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount must be positive")
	}
	return nil
}
