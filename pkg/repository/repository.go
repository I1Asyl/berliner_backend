package repository

import "xorm.io/xorm"

type Repository struct {
	Orm *xorm.Engine
}

func NewRepository() *Repository {
	return &Repository{Orm: SetupOrm()}
}
