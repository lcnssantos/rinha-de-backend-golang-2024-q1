package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/logging"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/rest"
	"github.com/mvmaasakkers/go-problemdetails"
	"gorm.io/gorm"
	"net/http"
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
			logging.Error(e.Request().Context(), err).Msg("error to bind payload")

			prob := problemdetails.New(
				http.StatusUnprocessableEntity,
				"UnprocessableEntity",
				"Unprocessable entity",
				err.Error(),
				"",
			)

			e.JSON(http.StatusUnprocessableEntity, prob)

			return err
		}

		customer, err := transactionService.Create(e.Request().Context(), payload.ID, payload.ToDomain())

		if err != nil {
			//logging.Error(e.Request().Context(), err).Msg("error to create transaction")

			if errors.Is(err, gorm.ErrRecordNotFound) {
				prob := problemdetails.New(
					http.StatusNotFound,
					"CustomerNotFound",
					"Customer not found",
					"",
					"",
				)

				e.JSON(http.StatusNotFound, prob)

				return err
			}

			if errors.Is(err, domain.ErrLimitExceeded) {
				prob := problemdetails.New(
					http.StatusUnprocessableEntity,
					"LimitExceeded",
					"Limit exceeded",
					"",
					"",
				)

				e.JSON(http.StatusUnprocessableEntity, prob)
				return err
			}

			prob := problemdetails.New(
				http.StatusInternalServerError,
				"InternalServerError",
				"Internal server error",
				"",
				"",
			)

			e.JSON(http.StatusInternalServerError, prob)
			return err
		}

		e.JSON(http.StatusOK, customerDto{}.FromDomain(customer))

		return nil
	}
}
