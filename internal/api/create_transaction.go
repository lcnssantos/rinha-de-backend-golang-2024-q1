package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/rest"
	"gorm.io/gorm"
)

type createTransactionDto struct {
	ID          uint64                 `param:"id" validate:"required,gt=0"`
	Amount      uint32                 `json:"valor" validate:"required,gt=0"`
	Type        domain.TransactionType `json:"tipo" validate:"required,oneof=c d"`
	Description string                 `json:"descricao" validate:"required,max=10,min=1"`
}

type customerDto struct {
	Limit  uint32 `json:"limite" validate:"required,gt=0"`
	Amount int32  `json:"saldo" validate:"required,gt=0"`
}

func (c customerDto) FromDomain(d domain.Customer) customerDto {
	return customerDto{
		Limit:  d.Limit,
		Amount: d.Amount,
	}
}

func (c createTransactionDto) ToDomain() domain.Transaction {
	return domain.Transaction{
		Amount:      c.Amount,
		Type:        c.Type,
		Description: c.Description,
		CustomerID:  c.ID,
	}
}

func createTransaction(transactionService domain.TransactionService) echo.HandlerFunc {
	return func(e echo.Context) error {
		payload, err := rest.Bind[createTransactionDto](e)

		if err != nil {
			e.NoContent(http.StatusUnprocessableEntity)

			return err
		}

		customer, err := transactionService.Create(e.Request().Context(), payload.ID, payload.ToDomain())

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				e.NoContent(http.StatusNotFound)

				return err
			}

			if errors.Is(err, domain.ErrLimitExceeded) {
				e.NoContent(http.StatusUnprocessableEntity)
				return err
			}

			e.NoContent(http.StatusInternalServerError)
			return err
		}

		e.JSON(http.StatusOK, customerDto{}.FromDomain(customer))

		return nil
	}
}
