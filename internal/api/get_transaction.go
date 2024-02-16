package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"

	"github.com/lcnssantos/rinha-de-backend/internal/domain"
)

type TransactionDto struct {
	Amount      uint32    `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

func (t TransactionDto) FromDomain(d domain.Transaction) TransactionDto {
	return TransactionDto{
		Amount:      d.Amount,
		Type:        string(d.Type),
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
	}
}

type CustomerAmount struct {
	Amount int32     `json:"total"`
	Date   time.Time `json:"data_extrato"`
	Limit  uint32    `json:"limite"`
}

func (c CustomerAmount) FromDomain(d domain.Customer) CustomerAmount {
	return CustomerAmount{
		Amount: d.Amount,
		Date:   time.Now(),
		Limit:  d.Limit,
	}
}

type Statement struct {
	Amount           CustomerAmount   `json:"saldo"`
	LastTransactions []TransactionDto `json:"ultimas_transacoes"`
}

func (a Statement) FromDomain(customer domain.Customer) Statement {
	t := []TransactionDto{}

	for _, v := range customer.Transactions {
		t = append(t, TransactionDto{}.FromDomain(v))
	}

	return Statement{
		Amount:           CustomerAmount{}.FromDomain(customer),
		LastTransactions: t,
	}
}

func getTransactions(service domain.TransactionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		id, err := c.ParamsInt("id")

		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		customer, err := service.GetTransactions(ctx, id)

		if err != nil {
			if errors.Is(err, domain.ErrCustomerNotFound) {
				return c.SendStatus(http.StatusNotFound)
			}

			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Status(http.StatusOK).JSON(Statement{}.FromDomain(*customer))
	}

}
