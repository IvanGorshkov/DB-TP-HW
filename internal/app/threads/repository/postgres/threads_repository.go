package repository

import (
	"database/sql"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"fmt"
)

type ThreadsRepository struct {
	dbConn *sql.DB
}

func NewThreadsRepository(conn *sql.DB) threads.ThreadsRepository {
	return &ThreadsRepository{
		dbConn: conn,
	}
}

 
func (tr *ThreadsRepository) ThreadById(id int) (*models.Thread, error) {
	var thread models.Thread

	err := tr.dbConn.QueryRow(`SELECT slug, id, forum from thread where id = $1`, id).Scan(&thread.Slug, &thread.Id, &thread.Forum)
	if err != nil {
		return nil, err
	}
	return &thread, nil

}

func (tr *ThreadsRepository) ThreadBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread
	err := tr.dbConn.QueryRow(`SELECT slug, id, forum  from thread where LOWER(slug) = LOWER($1)`, slug).Scan(&thread.Slug, &thread.Id, &thread.Forum)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (tr *ThreadsRepository) CreatePost(posts []*models.Post) ([]*models.Post, error) {
	query := `INSERT INTO posts (parent, author, message, forum, thread)
			VALUES `
	for i, post := range posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("(NULLIF(%d, 0), '%s', '%s', '%s', %d)", post.Parent, post.Author,
			post.Message, post.Forum, post.Thread)
	}
	query += " returning id, parent, author, message, is_edited, forum, thread, created"
	res, err := tr.dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	newPosts := make([]*models.Post, 0)
	var parent sql.NullInt64
	defer res.Close()
	for res.Next() {
		post := &models.Post{}
		err = res.Scan(&post.ID, &parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			fmt.Println(err)
		}
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			return nil, err
		}
		newPosts = append(newPosts, post)
	}

	return newPosts, nil
}