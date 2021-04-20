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
		where LOWER(nickname) = LOWER($1)
		) then TRUE else FALSE end`, forum.User).Scan(&isFind)
	if err != nil {
		return nil, err
	}

	if isFind == false {
		return nil, errors.New("404")
	}
	
	queryForum, err := fr.dbConn.Query(`SELECT title, nickname, slug, post, threads from forum where LOWER(slug) = LOWER($1)`,forum.Slug)
	if err != nil {
		return nil, err
	}

	defer queryForum.Close()

	for queryForum.Next() {
		var forum_409 models.Forum 
		err = queryForum.Scan(&forum_409.Title, &forum_409.User, &forum_409.Slug, &forum_409.Posts, &forum_409.Threads)
		if err != nil {
			return nil, err
		}
		return &forum_409, errors.New("409")
	}
	
	var user string;
	err = fr.dbConn.QueryRow(`SELECT nickname from users where LOWER(nickname) = LOWER($1)`,forum.User).Scan(&user)
	forum.User = user
	fmt.Println(err, user)
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
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (fr *ForumRepository) Detail(slug string) (*models.Forum, error) {
	var forum models.Forum
	err := fr.dbConn.QueryRow(`SELECT title, nickname, slug, post, threads from forum where LOWER(slug) = LOWER($1)`, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows 
		}
		return nil, err
	}

	return &forum, nil
}