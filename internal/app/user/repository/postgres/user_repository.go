package repository

import (
	"errors"
	"github.com/jackc/pgx"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
)

type UserRepository struct {
	dbConn *pgx.ConnPool
}

func NewUserRepository(conn *pgx.ConnPool) user.UserRepository {
	return &UserRepository{
		dbConn: conn,
	}
}

 
func(ur *UserRepository) Create(user *models.User) ([]*models.User, error) {
	tx, err := ur.dbConn.Begin()
	if err != nil {
		return nil, err
	}


	var users []*models.User
	queryUser, err := ur.dbConn.Query(`SELECT nickname, fullname, email, about FROM users WHERE nickname = $1 or email = $2`, user.Nickname, user.Email)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	defer queryUser.Close()

	for queryUser.Next() {
		var user_409 models.User 
		err = queryUser.Scan(&user_409.Nickname, &user_409.Fullname, &user_409.Email, &user_409.About)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
		users = append(users, &user_409)
	}
	if len(users) != 0 {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return users, errors.New("409")
	}

	_,_ = tx.Exec(`
		INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)
	`, user.Nickname, user.Fullname, user.Email, user.About)


	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	users = append(users, user)

	return users, nil
}

func(ur *UserRepository) GetProfile(nickname string) (*models.User, error) {

	var user models.User 
	err := ur.dbConn.QueryRow(`SELECT nickname, fullname, email, about FROM users WHERE nickname = $1`, nickname).Scan(
		&user.Nickname, &user.Fullname, &user.Email, &user.About,
	)
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func(ur *UserRepository) UpdateProfile(user *models.User) (*models.User, error) {
	var newUser models.User

	tx, err := ur.dbConn.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		`UPDATE users SET 	  email=COALESCE(NULLIF($1, ''), email), 
							  about=COALESCE(NULLIF($2, ''), about), 
							  fullname=COALESCE(NULLIF($3, ''), fullname) 
							  WHERE nickname = $4 
							  RETURNING nickname, fullname, about, email`,
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

		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}


	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	
	return &newUser, nil
}