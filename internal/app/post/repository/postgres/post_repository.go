package repository

import (
	"database/sql"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"time"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/post"
)

type PostRepository struct {
	dbConn *pgx.ConnPool
}

func NewPostRepository(conn *pgx.ConnPool) post.PostRepository {
	return &PostRepository{
		dbConn: conn,
	}
}


func (pr *PostRepository) Update(id int, post models.Post) (*models.Post, error) {
	var newPost models.Post

	var created time.Time
	var parent sql.NullInt64
	err := pr.dbConn.QueryRow(
		`UPDATE posts 
		 SET message=COALESCE(NULLIF($1, ''), message),
		 is_edited = CASE WHEN $1 = '' OR message = $1 THEN is_edited ELSE true END
	  	 WHERE id=$2 
		 RETURNING parent, author, message, is_edited, forum, thread, created`,
							 post.Message,
							 id,
	).Scan(
		&parent,
		&newPost.Author,
		&newPost.Message,
		&newPost.IsEdited,
		&newPost.Forum,
		&newPost.Thread,
		&created,
	)

	newPost.Created = strfmt.DateTime(created.UTC()).String()

	newPost.ID = id
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

	var created time.Time
	var parent sql.NullInt64
	err := pr.dbConn.QueryRow(`SELECT parent, author, message, is_edited, forum, thread, created 
								from posts 
								WHERE id = $1`, id).Scan(&parent ,&post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &created)
	
	post.ID = id
	if parent.Valid {
		post.Parent = int(parent.Int64)
	}
	post.Created = strfmt.DateTime(created.UTC()).String()

	if err != nil {
		return nil, err
	}
	
	return &post, nil
}