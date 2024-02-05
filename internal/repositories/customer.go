package repositories

import (
	"context"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"gorm.io/gorm"
)

type customerImpl struct {
	db *gorm.DB
}

func (c customerImpl) SelectForUpdate(ctx context.Context, ID uint64) (bool, error) {
	var id uint64

	return id > 0, c.db.WithContext(ctx).Raw("SELECT id as exists FROM customers WHERE id = ? FOR UPDATE", ID).Scan(&id).Error
}

func (c customerImpl) AddAmount(ctx context.Context, ID uint64, amount int32) error {
	return c.db.WithContext(ctx).Model(&domain.Customer{}).Where("id = ?", ID).Update("amount", gorm.Expr("amount + ?", amount)).Error
}

func (c customerImpl) FindOne(ctx context.Context, id uint64) (domain.Customer, error) {
	var customer domain.Customer
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&customer).Error
	return customer, err
}

func (c customerImpl) BeginTransaction() *gorm.DB {
	return c.db.Begin()
}

func (c customerImpl) WithTransaction(db *gorm.DB) domain.CustomerRepository {
	c.db = db
	return c
}

func NewCustomer(db *gorm.DB) domain.CustomerRepository {
	return &customerImpl{
		db: db,
	}
}
