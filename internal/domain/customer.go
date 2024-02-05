package domain

import (
	"errors"

	"gorm.io/gorm"
)

type Customer struct {
	ID           uint64        `json:"id" gorm:"primaryKey"`
	Limit        uint32        `json:"limit"`
	Amount       int32         `json:"amount"`
	Transactions []Transaction `json:"transactions"`
}

var ErrCustomerNotFound = gorm.ErrRecordNotFound
var ErrLimitExceeded = errors.New("limit exceeded")
