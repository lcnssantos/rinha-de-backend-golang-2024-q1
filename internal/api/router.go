package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/postgres"
)

func RoutesFactory(transactionService domain.TransactionService, postgres postgres.Postgres) func(f fiber.Router) {
	return func(f fiber.Router) {
		f.Post("/clientes/:id/transacoes", createTransaction(transactionService))
		f.Get("/clientes/:id/extrato", getTransactions(transactionService))

		f.Delete("", func(c *fiber.Ctx) error {
			err := postgres.DB().Exec("TRUNCATE TABLE transactions RESTART IDENTITY CASCADE").Error

			if err != nil {
				return err
			}

			err = postgres.DB().Exec("UPDATE customers SET amount = 0 WHERE true").Error

			if err != nil {
				return err
			}

			return nil
		})
	}
}
