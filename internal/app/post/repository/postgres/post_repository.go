package repository

import (
	"database/sql"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/post"
)

type PostRepository struct {
	dbConn *sql.DB
}

func NewPostRepository(conn *sql.DB) post.PostRepository {
	return &PostRepository{
		dbConn: conn,
	}
}


func (pr *PostRepository) Update(id int, post models.Post) (*models.Post, error) {
	var newPost models.Post

	var parent sql.NullInt64
	err := pr.dbConn.QueryRow(
		`UPDATE posts SET message=COALESCE(NULLIF($1, ''), message),
		is_edited = CASE WHEN $1 = '' OR message = $1 THEN is_edited ELSE true END
							 WHERE id=$2 
							 RETURNING id, parent, author, message, is_edited, forum, thread, created`,
							 post.Message,
							 id,
	).Scan(
		&newPost.ID,
		&parent,
		&newPost.Author,
		&newPost.Message,
		&newPost.IsEdited,
		&newPost.Forum,
		&newPost.Thread,
		&newPost.Created,
	)

	if parent.Valid {
		post.Parent = int(parent.Int64)
	}

	if err != nil {
		return nil, err
	}

	return &newPost, nil
}


func (pr *PostRepository) GetPostById(id int) (*models.Post, error) {
	var post models.Post

	var parent sql.NullInt64
	err := pr.dbConn.QueryRow(`SELECT id, parent, author, message, is_edited, forum, thread, created 
								from posts 
								WHERE id = $1`, id).Scan(&post.ID, &parent ,&post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	

	if parent.Valid {
		post.Parent = int(parent.Int64)
	}

	if err != nil {
		return nil, err
	}
	
	return &post, nil
}