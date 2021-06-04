package repository

import (
	"database/sql"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"strings"
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

func (pr *PostRepository) GetPostFullbyId(id int, related string) (*models.PostFull, error) {

	var parent sql.NullInt64
	var created time.Time
	fullPost := &models.PostFull{
		Post: models.Post{},
	}

	if strings.Contains(related, "forum") && strings.Contains(related, "thread") && strings.Contains(related, "user") {
		var created2 time.Time
		fullPost.Author = &models.User{}
		fullPost.Thread = &models.Thread{}
		fullPost.Forum = &models.Forum{}
		err := pr.dbConn.QueryRow(`
					SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
					   th.slug,
					   th.id,
					   th.forum,
					   th.title,
					   th.author,
					   th.message,
					   th.votes,
					   th.created,
					   f.title,
					   f.nickname,
					   f.slug,
					   f.post,
					   f.threads,
			   		   u.nickname,
			   		   u.fullname,
			   		   u.about,
			   		   u.email
				FROM posts as p
				INNER JOIN thread th on th.id = p.thread
				INNER JOIN forum f on f.slug = p.forum
				INNER JOIN users u on u.nickname = p.author
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Id,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&created2,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads,
			&fullPost.Author.Nickname,
			&fullPost.Author.Fullname,
			&fullPost.Author.About,
			&fullPost.Author.Email,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		fullPost.Thread.Created = strfmt.DateTime(created2.UTC()).String()
		return fullPost, err
	}

	if strings.Contains(related, "forum") && strings.Contains(related, "thread")  {
		var created2 time.Time
		fullPost.Thread = &models.Thread{}
		fullPost.Forum = &models.Forum{}
		err := pr.dbConn.QueryRow(`
					SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
					   th.slug,
					   th.id,
					   th.forum,
					   th.title,
					   th.author,
					   th.message,
					   th.votes,
					   th.created,
					   f.title,
					   f.nickname,
					   f.slug,
					   f.post,
					   f.threads
				FROM posts as p
				INNER JOIN thread th on th.id = p.thread
				INNER JOIN forum f on f.slug = p.forum
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Id,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&created2,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		fullPost.Thread.Created = strfmt.DateTime(created2.UTC()).String()
		return fullPost, err
	}

	if strings.Contains(related, "forum") && strings.Contains(related, "user") {
		fullPost.Forum = &models.Forum{}
		fullPost.Author = &models.User{}
		err := pr.dbConn.QueryRow(`
					SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
					   f.title,
					   f.nickname,
					   f.slug,
					   f.post,
					   f.threads,
			   		   u.nickname,
			   		   u.fullname,
			   		   u.about,
			   		   u.email
				FROM posts as p
				INNER JOIN forum f on f.slug = p.forum
				INNER JOIN users u on u.nickname = p.author
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads,
			&fullPost.Author.Nickname,
			&fullPost.Author.Fullname,
			&fullPost.Author.About,
			&fullPost.Author.Email,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		return fullPost, err
	}

	if strings.Contains(related, "thread") && strings.Contains(related, "user") {
		var created2 time.Time
		fullPost.Thread = &models.Thread{}
		fullPost.Author = &models.User{}
		err := pr.dbConn.QueryRow(`
					SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
					   th.slug,
					   th.id,
					   th.forum,
					   th.title,
					   th.author,
					   th.message,
					   th.votes,
					   th.created,
			   		   u.nickname,
			   		   u.fullname,
			   		   u.about,
			   		   u.email
				FROM posts as p
				INNER JOIN thread th on th.id = p.thread
				INNER JOIN users u on u.nickname = p.author
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Id,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&created2,
			&fullPost.Author.Nickname,
			&fullPost.Author.Fullname,
			&fullPost.Author.About,
			&fullPost.Author.Email,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		fullPost.Thread.Created = strfmt.DateTime(created2.UTC()).String()
		return fullPost, err
	}

	if strings.Contains(related, "user") {
		fullPost.Author = &models.User{}
		err := pr.dbConn.QueryRow(`
				SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
			   		   u.nickname,
			   		   u.fullname,
			   		   u.about,
			   		   u.email
				FROM posts as p
				INNER JOIN users u on u.nickname = p.author
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Author.Nickname,
			&fullPost.Author.Fullname,
			&fullPost.Author.About,
			&fullPost.Author.Email,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		return fullPost, err
	}

	if strings.Contains(related, "thread") {
		var created2 time.Time
		fullPost.Thread = &models.Thread{}
		err := pr.dbConn.QueryRow(`
					SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
					   th.slug,
					   th.id,
					   th.forum,
					   th.title,
					   th.author,
					   th.message,
					   th.votes,
					   th.created
				FROM posts as p
				INNER JOIN thread th on th.id = p.thread
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Thread.Slug,
			&fullPost.Thread.Id,
			&fullPost.Thread.Forum,
			&fullPost.Thread.Title,
			&fullPost.Thread.Author,
			&fullPost.Thread.Message,
			&fullPost.Thread.Votes,
			&created2,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		fullPost.Thread.Created = strfmt.DateTime(created2.UTC()).String()
		return fullPost, err
	}

	if strings.Contains(related, "forum") {
		fullPost.Forum = &models.Forum{}
		err := pr.dbConn.QueryRow(`
					SELECT 
					   p.forum,
  					   p.thread,
			   		   p.author,
			   		   p.message,
			   		   p.is_edited,
			   		   p.created,
			   		   p.parent,
					   f.title,
					   f.nickname,
					   f.slug,
					   f.post,
					   f.threads
				FROM posts as p
				INNER JOIN forum f on f.slug = p.forum
				WHERE p.id = $1
				LIMIT 1
		`,
			id).Scan(
			&fullPost.Post.Forum,
			&fullPost.Post.Thread,
			&fullPost.Post.Author,
			&fullPost.Post.Message,
			&fullPost.Post.IsEdited,
			&created,
			&parent,
			&fullPost.Forum.Title,
			&fullPost.Forum.User,
			&fullPost.Forum.Slug,
			&fullPost.Forum.Posts,
			&fullPost.Forum.Threads,
		)
		fullPost.Post.ID = id
		if err != nil {
			return nil, err
		}
		if parent.Valid {
			fullPost.Post.Parent = int(parent.Int64)
		}
		fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()
		return fullPost, err
	}

	err := pr.dbConn.QueryRow(`SELECT 
									parent, 
		    						author, 
		    						message, 
		    						is_edited,
		    						forum,
		    						thread,
		    						created 
								from posts 
								WHERE id = $1`, id).Scan(
									&parent ,&fullPost.Post.Author, &fullPost.Post.Message, &fullPost.Post.IsEdited, &fullPost.Post.Forum, &fullPost.Post.Thread, &created)

	fullPost.Post.ID = id
	if parent.Valid {
		fullPost.Post.Parent = int(parent.Int64)
	}
	fullPost.Post.Created = strfmt.DateTime(created.UTC()).String()

	if err != nil {
		return nil, err
	}


	return fullPost, nil
}