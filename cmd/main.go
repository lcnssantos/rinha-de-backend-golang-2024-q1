package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lcnssantos/rinha-de-backend/internal/api"
	"github.com/lcnssantos/rinha-de-backend/internal/env"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/environment"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/logging"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/postgres"
	"github.com/lcnssantos/rinha-de-backend/internal/lib/rest"
	"github.com/lcnssantos/rinha-de-backend/internal/repositories"
	"github.com/lcnssantos/rinha-de-backend/internal/services"
	"os"
)

func main() {
	ctx := context.TODO()

	logging.Init(os.Getenv("LOG_LEVEL"))

	environmentVariables, err := environment.LoadEnvironment[env.Environment]()

	if err != nil {
		logging.Panic(ctx, err).Msg("error to load environment variables")
		return
	}

	pg := postgres.New(
		postgres.NewConfig(
			environmentVariables.DatabaseHost,
			environmentVariables.DatabasePort,
			environmentVariables.DatabaseUser,
			environmentVariables.DatabasePass,
			environmentVariables.DatabaseName,
			environmentVariables.DatabaseSSLMode,
		),
	)

	err = pg.Connect()

	if err != nil {
		logging.Panic(ctx, err).Msg("error to connect to database")
		return
	}

	migrationsPath := fmt.Sprintf("file://%s", "./migrations")

	err = pg.MigrateUp(migrationsPath)

	if err != nil {
		logging.Panic(ctx, err).Msg("error to migrate database")
		return
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = &rest.CustomValidator{
		Validator: validator.New(),
	}

	transactionRepository := repositories.NewTransaction(pg.DB().Debug())
	customerRepository := repositories.NewCustomer(pg.DB().Debug())

	transactionService := services.NewTransactionService(transactionRepository, customerRepository)

	api.RoutesFactory(transactionService)(e.Group(""))

	routes := e.Routes()

	for _, route := range routes {
		logging.Info(ctx).
			Str("method", route.Method).
			Str("path", route.Path).
			Msg("route settled")
	}

	err = e.Start(fmt.Sprintf(":%s", environmentVariables.Port))

	if err != nil {
		logging.Panic(ctx, err).Msg("error to start http server")
		return
	}
}
