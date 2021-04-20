package user

import "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"


type UserRepository interface {
	Create(user *models.User) ([]*models.User, error)
	GetProfile(nickname string) (*models.User, error)
	UpdateProfile(user *models.User) (*models.User, error)
}
