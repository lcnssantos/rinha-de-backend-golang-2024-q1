package services

import (
	"context"
	"strconv"

	"gorm.io/gorm/clause"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"gorm.io/gorm"
)

type transactionService struct {
	gorm *gorm.DB
}

func (t transactionService) GetTransactions(ctx context.Context, id string) (domain.Customer, error) {
	_id, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return domain.Customer{}, err
	}

	var customer domain.Customer

	err = t.gorm.WithContext(ctx).Preload("Transactions", func(tx *gorm.DB) *gorm.DB {
		return tx.Limit(10).Order("created_at DESC")
	}).Where("id = ?", _id).First(&customer).Error

	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}

func (t transactionService) Create(ctx context.Context, id uint64, transaction domain.Transaction) (domain.Customer, error) {
	var customer domain.Customer

	err := t.gorm.Transaction(func(tx *gorm.DB) error {
		expression := "amount + ?"

		if transaction.Type == domain.TransactionTypeDebit {
			expression = "amount - ?"
		}

		err := tx.WithContext(ctx).
			Model(&customer).
			Clauses(clause.Returning{
				Columns: []clause.Column{
					{Name: "limit"},
					{
						Name: "amount",
					},
				},
			}).
			Where("id = ?", id).
			Update("amount", gorm.Expr(expression, transaction.Amount)).Error

		if err != nil {
			pgErr, ok := err.(*pgconn.PgError)

			if ok && pgErr.ConstraintName == "customer_limit_check" {
				return domain.ErrLimitExceeded
			}

			return err
		}

		err = tx.WithContext(ctx).Create(&transaction).Error

		if err != nil {
			pgErr, ok := err.(*pgconn.PgError)

			if ok && pgErr.ConstraintName == "fk_customer" {
				return domain.ErrCustomerNotFound
			}

			return err
		}

		return nil
	})

	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}

func NewTransactionService(gorm *gorm.DB) domain.TransactionService {
	return transactionService{
		gorm,
	}
}
