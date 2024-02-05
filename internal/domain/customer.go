package domain

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Customer struct {
	ID     uint64 `json:"id" gorm:"primaryKey"`
	Limit  uint32 `json:"limit"`
	Amount int32  `json:"amount"`
}

type CustomerRepository interface {
	BeginTransaction() *gorm.DB
	WithTransaction(db *gorm.DB) CustomerRepository
	SelectForUpdate(ctx context.Context, ID uint64) (bool, error)
	AddAmount(ctx context.Context, ID uint64, amount int32) error
	FindOne(ctx context.Context, id uint64) (Customer, error)
}

var ErrCustomerNotFound = gorm.ErrRecordNotFound
var ErrLimitExceeded = errors.New("limit exceeded")
