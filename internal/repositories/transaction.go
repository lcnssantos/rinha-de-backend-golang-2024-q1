package repositories

import (
	"context"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"gorm.io/gorm"
)

type transactionImpl struct {
	db *gorm.DB
}

func (t transactionImpl) FindAll(ctx context.Context, id uint64) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	return transactions, t.db.WithContext(ctx).Where("customer_id = ?", id).Find(&transactions).Error
}

func (t transactionImpl) BeginTransaction(ctx context.Context) *gorm.DB {
	return t.db.WithContext(ctx).Begin()
}

func (t transactionImpl) WithTransaction(db *gorm.DB) domain.TransactionRepository {
	t.db = db
	return t
}

func (t transactionImpl) Create(ctx context.Context, transaction domain.Transaction) error {
	return t.db.WithContext(ctx).Create(&transaction).Error
}

func NewTransaction(db *gorm.DB) domain.TransactionRepository {
	return &transactionImpl{
		db: db,
	}
}
