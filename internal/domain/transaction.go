package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "c"
	TransactionTypeDebit  TransactionType = "d"
)

type Transaction struct {
	ID          uint64          `json:"id" gorm:"primaryKey"`
	Amount      uint32          `json:"amount"`
	CustomerID  int             `json:"customer_id" gorm:"foreignKey:CustomerID"`
	Description string          `json:"description"`
	Type        TransactionType `json:"type"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (t Transaction) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	return nil
}

type TransactionService interface {
	Create(context.Context, int, Transaction) (*Customer, error)
	GetTransactions(ctx context.Context, id int) (*Customer, error)
}
