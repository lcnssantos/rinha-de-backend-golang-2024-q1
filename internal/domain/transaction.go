package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "c"
	TransactionTypeDebit  TransactionType = "d"
)

type Transaction struct {
	ID          uint64          `json:"id" gorm:"primaryKey"`
	Amount      uint32          `json:"amount"`
	CustomerID  uint64          `json:"customer_id" gorm:"foreignKey:CustomerID"`
	Description string          `json:"description"`
	Type        TransactionType `json:"type"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (t Transaction) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	return nil
}

type TransactionRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	WithTransaction(db *gorm.DB) TransactionRepository
	Create(context.Context, Transaction) error
	FindAll(ctx context.Context, id uint64) ([]Transaction, error)
}

type TransactionService interface {
	Create(context.Context, uint64, Transaction) (Customer, error)
	GetTransactions(ctx context.Context, id string) (Customer, []Transaction, error)
}
