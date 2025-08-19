package repository

import (
	"Book_Exp/wire/repository/dao"
)

type Repository struct {
	dao dao.UserDao
}

func NewRepository(dao dao.UserDao) *Repository {
	return &Repository{
		dao: dao,
	}
}
