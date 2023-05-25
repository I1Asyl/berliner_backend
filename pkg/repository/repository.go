package repository

import "xorm.io/xorm"

type Repository struct {
	Orm *xorm.Engine
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

func NewRepository(orm *xorm.Engine) *Repository {
	return &Repository{Orm: orm}
}
