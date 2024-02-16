package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"

	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/rest"
	"gorm.io/gorm"
)

type createTransactionDto struct {
	//ID          uint64                 `param:"id" validate:"required,gt=0"`
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

func (c createTransactionDto) ToDomain(id int) domain.Transaction {
	return domain.Transaction{
		Amount:      c.Amount,
		Type:        c.Type,
		Description: c.Description,
		CustomerID:  id,
	}
}

func createTransaction(transactionService domain.TransactionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := c.ParamsInt("id")

		if err != nil {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}

		payload, err := rest.Bind[createTransactionDto](c)

		if err != nil {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}

		customer, err := transactionService.Create(ctx, id, payload.ToDomain(id))

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.SendStatus(http.StatusNotFound)
			}

			if errors.Is(err, domain.ErrLimitExceeded) {
				return c.SendStatus(http.StatusUnprocessableEntity)
			}

			return c.SendStatus(http.StatusInternalServerError)
		}

		c.Status(http.StatusOK).JSON(customerDto{}.FromDomain(*customer))

		return nil
	}
}
