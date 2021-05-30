package usecase

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"

	"github.com/jackc/pgx"
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
		if err == pgx.ErrNoRows {
			return nil, errors.NotFoundBody("Can't find user with nickname " + nickname + "\n")
		}

		return nil, errors.UnexpectedInternal(err)
	}
	return res, nil
}

func(us *UserUsecase) UpdateProfile(user *models.User) (*models.User, *errors.Error) {
	res, err := us.userRepo.UpdateProfile(user)

	if err == pgx.ErrNoRows {
		return nil, errors.NotFoundBody("Can't find user with nickname " + user.Nickname + "\n")
	}

	if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
		return nil, errors.ConflictErrorBody("This email is already registered by user: " + err.Error()+ "\n")
	}
	return res, nil
}