package usecase

import (
	"database/sql"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func(us *UserUsecase) Create(user *models.User) ([]*models.User, *errors.Error) {
	res, err := us.userRepo.Create(user)
	if err != nil {
		if err.Error() == "409" {
			return res, errors.CustomErrors[errors.ConflictError]
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}

func(us *UserUsecase) GetProfile(nickname string) (*models.User, *errors.Error) {
	res, err := us.userRepo.GetProfile(nickname)
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	if res == nil {
		return nil, errors.NotFoundBody("Can't find user with nickname " + nickname + "\n")
	}
	return res, nil
}

func(us *UserUsecase) UpdateProfile(user *models.User) (*models.User, *errors.Error) {
	res, err := us.userRepo.UpdateProfile(user)
	if err == sql.ErrNoRows {
		return nil, errors.NotFoundBody("Can't find user with nickname " + user.Nickname + "\n")
	}
	if err != nil {
		return nil, errors.ConflictErrorBody("Can't find user with nickname " + user.Nickname + "\n")
	}
	return res, nil
}