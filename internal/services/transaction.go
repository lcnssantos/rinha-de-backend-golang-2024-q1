package services

import (
	"context"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"strconv"
)

type transactionService struct {
	transactionRepository domain.TransactionRepository
	customerRepository    domain.CustomerRepository
}

func (t transactionService) GetTransactions(ctx context.Context, id string) (domain.Customer, []domain.Transaction, error) {
	_id, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return domain.Customer{}, []domain.Transaction{}, err
	}

	customer, err := t.customerRepository.FindOne(ctx, _id)

	if err != nil {
		return domain.Customer{}, []domain.Transaction{}, err
	}

	transactions, err := t.transactionRepository.FindAll(ctx, _id)

	if err != nil {
		return domain.Customer{}, []domain.Transaction{}, err
	}

	return customer, transactions, nil
}

func (t transactionService) Create(ctx context.Context, id uint64, transaction domain.Transaction) (domain.Customer, error) {
	tx := t.transactionRepository.BeginTransaction(ctx)

	defer tx.Rollback()

	exists, err := t.customerRepository.WithTransaction(tx).SelectForUpdate(ctx, id)

	if err != nil {
		return domain.Customer{}, err
	}

	if !exists {
		return domain.Customer{}, domain.ErrCustomerNotFound
	}

	err = t.transactionRepository.WithTransaction(tx).Create(ctx, transaction)

	if err != nil {
		return domain.Customer{}, err
	}

	var addedAmount int32

	if transaction.Type == domain.TransactionTypeCredit {
		addedAmount = int32(transaction.Amount)
	} else {
		addedAmount = -int32(transaction.Amount)
	}

	err = t.customerRepository.WithTransaction(tx).AddAmount(ctx, id, addedAmount)

	if err != nil {
		return domain.Customer{}, err
	}

	customer, err := t.customerRepository.WithTransaction(tx).FindOne(ctx, id)

	if err != nil {
		return domain.Customer{}, err
	}

	err = tx.Commit().Error

	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}

func NewTransactionService(transactionRepository domain.TransactionRepository, customerRepository domain.CustomerRepository) domain.TransactionService {
	return transactionService{
		transactionRepository,
		customerRepository,
	}
}
