package repository

import (
	"database/sql"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
	"context"
)

type UserRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserRepository {
	return &UserRepository{
		dbConn: conn,
	}
}

 
func(ur *UserRepository) Create(user *models.User) (*models.User, error) {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	query := tx.QueryRow(`
		INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)
	`, user.Nickname, user.Fullname, user.Email, user.About)

	err = query.Scan()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		
		queryUser, err := ur.dbConn.Query(`SELECT nickname, fullname, email, about FROM users WHERE nickname = $1`, user.Nickname)
		if err != nil {
			return nil, err
		}

		defer queryUser.Close()

		for queryUser.Next() {
			err = queryUser.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
			if err != nil {
				return nil, err
			}
		}
		return user, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func(us *UserRepository) GetProfile(nickname string) (*models.User, error) {
	return nil, nil
}

func(us *UserRepository) UpdateProfile(user *models.User) (*models.User, error) {
	return nil, nil
}