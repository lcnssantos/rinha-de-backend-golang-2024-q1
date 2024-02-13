package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

func getTransactions(service domain.TransactionService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		customer, err := service.GetTransactions(c.Request().Context(), id)

		if err != nil {
			if errors.Is(err, domain.ErrCustomerNotFound) {
				return c.NoContent(http.StatusNotFound)
			}

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, Statement{}.FromDomain(customer))
	}
}
