package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lcnssantos/rinha-de-backend/internal/domain"
)

func RoutesFactory(transactionService domain.TransactionService) func(g *echo.Group) {
	return func(g *echo.Group) {
		g.POST("/clientes/:id/transacoes", createTransaction(transactionService))
		g.GET("/clientes/:id/extrato", getTransactions(transactionService))
	}
}
