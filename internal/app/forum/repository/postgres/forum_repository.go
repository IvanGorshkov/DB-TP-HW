package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/forum"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type ForumRepository struct {
	dbConn *sql.DB
}

func NewForumRepository(conn *sql.DB) forum.ForumRepository {
	return &ForumRepository{
		dbConn: conn,
	}
}

 
func(fr *ForumRepository) Create(forum *models.Forum) (*models.Forum, error){
	var isFind = false
	err := fr.dbConn.QueryRow(`
	select case when EXISTS (
		select 1 
		from users
		where nickname = $1
		) then TRUE else FALSE end`, forum.User).Scan(&isFind)
	if err != nil {
		return nil, err
	}

	if isFind == false {
		return nil, errors.New("404")
	}
	
	tx, err := fr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	
	query := tx.QueryRow(`
		INSERT INTO forum (title, nickname, slug) VALUES ($1, $2, $3) returning id
	`, forum.Title, forum.User, forum.Slug)

	id := 0
	err = query.Scan(&id)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		
		err := fr.dbConn.QueryRow(`SELECT title, nickname, slug, post, threads from forum where slug = $1`, forum.Slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
		if err != nil {
			return nil, err
		}
		return forum, errors.New("409")
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (fr *ForumRepository) Detail(slug string) (*models.Forum, error) {
	var forum models.Forum
	err := fr.dbConn.QueryRow(`SELECT title, nickname, slug, post, threads from forum where slug = $1`, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows 
		}
		return nil, err
	}

	return &forum, nil
}