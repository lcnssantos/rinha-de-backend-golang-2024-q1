package postgres

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres interface {
	MigrateUp(path string) error
	MigrateDown(path string) error
	Connect() error
	DB() *gorm.DB
	WithPoolConfig(poolConfig poolConfig) Postgres
}

type postgresClientImpl struct {
	config     config
	poolConfig *poolConfig
	db         *gorm.DB
}

func New(config config) Postgres {
	return &postgresClientImpl{
		config: config,
	}
}

func (p *postgresClientImpl) WithPoolConfig(poolConfig poolConfig) Postgres {
	p.poolConfig = &poolConfig

	return p
}

func (p *postgresClientImpl) MigrateUp(path string) error {
	sqlInstance, err := p.db.DB()

	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlInstance, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		p.config.Database,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (p *postgresClientImpl) MigrateDown(path string) error {
	sqlInstance, err := p.db.DB()

	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlInstance, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		p.config.Database,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Down()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (p *postgresClientImpl) DB() *gorm.DB {
	return p.db
}

func (p *postgresClientImpl) Connect() error {
	db, err := gorm.Open(gormpostgres.Open(p.config.string()), &gorm.Config{})

	if err != nil {
		return err
	}

	p.db = db

	if p.poolConfig != nil {
		sqlDB, err := p.db.DB()

		if err != nil {
			return err
		}

		if p.poolConfig.maxIdle < 1 {
			return ErrInvalidPoolConfiguration
		}

		if p.poolConfig.maxOpen < 1 {
			return ErrInvalidPoolConfiguration
		}

		if p.poolConfig.maxLifeTime < 1 {
			return ErrInvalidPoolConfiguration
		}

		sqlDB.SetMaxIdleConns(p.poolConfig.maxIdle)
		sqlDB.SetMaxOpenConns(p.poolConfig.maxOpen)
		sqlDB.SetConnMaxLifetime(p.poolConfig.maxLifeTime)
	}

	return nil
}
