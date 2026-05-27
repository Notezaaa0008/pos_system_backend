package users

import (
	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUserRepository (db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

