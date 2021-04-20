package user

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type UserUsecase interface {
	Create(user *models.User) ([]*models.User, *errors.Error)
	GetProfile(nickname string) (*models.User, *errors.Error)
	UpdateProfile(user *models.User) (*models.User, *errors.Error)
}
