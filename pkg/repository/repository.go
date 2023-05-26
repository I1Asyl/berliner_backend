package repository

import (
	"github.com/golang-migrate/migrate/v4"
	"xorm.io/xorm"
)

type Repository struct {
	Orm       *xorm.Engine
	Migration *migrate.Migrate
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

func NewRepository(dsn string, migrationsPath string) *Repository {
	return &Repository{Orm: NewOrm(dsn), Migration: NewMigration(dsn, migrationsPath)}
}
