package repository

import (
	"database/sql"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
)

type UserRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserRepository {
	return &UserRepository{
		dbConn: conn,
	}
}

 
func(us *UserRepository) Create(user *models.User) (*models.User, error) {
	return nil, nil
}

func(us *UserRepository) GetProfile(nickname string) (*models.User, error) {
	return nil, nil
}

func(us *UserRepository) UpdateProfile(user *models.User) (*models.User, error) {
	return nil, nil
}