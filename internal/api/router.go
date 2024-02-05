package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/postgres"
)

func RoutesFactory(transactionService domain.TransactionService, postgres postgres.Postgres) func(g *echo.Group) {
	return func(g *echo.Group) {
		//g.Use(loggingMiddleware)
		g.POST("/clientes/:id/transacoes", createTransaction(transactionService))
		g.GET("/clientes/:id/extrato", getTransactions(transactionService))
		g.DELETE("", func(c echo.Context) error {
			err := postgres.DB().Debug().Exec("TRUNCATE TABLE transactions RESTART IDENTITY CASCADE").Error

			if err != nil {
				return err
			}

			err = postgres.DB().Debug().Exec("UPDATE customers SET amount = 0 WHERE true").Error

			if err != nil {
				return err
			}
			return nil
		})
	}
}
