package repository

import (
	"context"
	"database/sql"
	"errors"

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

 
func(ur *UserRepository) Create(user *models.User) ([]*models.User, error) {
	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}


	var users []*models.User
	queryUser, err := ur.dbConn.Query(`SELECT nickname, fullname, email, about FROM users WHERE LOWER(nickname) = LOWER($1) or LOWER(email) = LOWER($2)`, user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}

	defer queryUser.Close()

	for queryUser.Next() {
		var user_409 models.User 
		err = queryUser.Scan(&user_409.Nickname, &user_409.Fullname, &user_409.Email, &user_409.About)
		if err != nil {
			return nil, err
		}
		users = append(users, &user_409)
	}
	if len(users) != 0 {
		return users, errors.New("409")
	}

	query := tx.QueryRow(`
		INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4) returning id
	`, user.Nickname, user.Fullname, user.Email, user.About)

	id := 0
	err = query.Scan(&id)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	users = append(users, user)

	return users, nil
}

func(ur *UserRepository) GetProfile(nickname string) (*models.User, error) {
	queryUser, err := ur.dbConn.Query(`SELECT nickname, fullname, email, about FROM users WHERE LOWER(nickname)= LOWER($1)`, nickname)
		if err != nil {
			return nil, err
		}

		defer queryUser.Close()

		for queryUser.Next() {
			var user models.User 
			err = queryUser.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
			if err != nil {
				return nil, err
			}
			return &user, nil
		}
		return nil, nil
}

func(ur *UserRepository) UpdateProfile(user *models.User) (*models.User, error) {
	var newUser models.User

	tx, err := ur.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	var findUser models.User
	err = tx.QueryRow(
		`SELECT nickname from users where LOWER(email) = LOWER($1)`,
		user.Email,
	).Scan(&findUser.Nickname)

	if findUser.Nickname != "" {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New(findUser.Nickname)
	}

	err = tx.QueryRow(
		`UPDATE users SET email=COALESCE(NULLIF($1, ''), email), 
							  about=COALESCE(NULLIF($2, ''), about), 
							  fullname=COALESCE(NULLIF($3, ''), fullname) WHERE LOWER(nickname)=LOWER($4) RETURNING nickname, fullname, about, email`,
		user.Email,
		user.About,
		user.Fullname,
		user.Nickname,
	).Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}

			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows 
			}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	
	return &newUser, nil
}