package usecase

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewProductUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func(us *UserUsecase) Create(user *models.User) (*models.User, error) {
	return nil, nil
}

func(us *UserUsecase) GetProfile(nickname string) (*models.User, error) {
	return nil, nil
}

func(us *UserUsecase) UpdateProfile(user *models.User) (*models.User, error) {
	return nil, nil
}