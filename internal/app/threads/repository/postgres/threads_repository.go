package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
)

type ThreadsRepository struct {
	dbConn *sql.DB
}

func NewThreadsRepository(conn *sql.DB) threads.ThreadsRepository {
	return &ThreadsRepository{
		dbConn: conn,
	}
}

func (tr *ThreadsRepository) UpdateById(thread *models.Thread, id int) (*models.Thread, error) {
	var new_thread models.Thread
	query := `UPDATE thread SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) WHERE id = $3 RETURNING *`
	tr.dbConn.QueryRow(query, thread.Title, thread.Message, id).Scan(
		&new_thread.Id,
		&new_thread.Title,
		&new_thread.Author,
		&new_thread.Forum,
		&new_thread.Message,
		&new_thread.Votes,
		&new_thread.Created,
		&new_thread.Slug,
	)

	return &new_thread, nil
}

func (tr *ThreadsRepository) UpdateBySlug(thread *models.Thread, slug string) (*models.Thread, error) {
	var new_thread models.Thread
	query := `UPDATE thread SET title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message) WHERE LOWER(slug) = LOWER($3) RETURNING *`
	tr.dbConn.QueryRow(query, thread.Title, thread.Message, slug).Scan(
		&new_thread.Id,
		&new_thread.Title,
		&new_thread.Author,
		&new_thread.Forum,
		&new_thread.Message,
		&new_thread.Votes,
		&new_thread.Created,
		&new_thread.Slug,
	)

	return &new_thread, nil
}

func FormQueryFlatSort(limit, threadID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE thread = $1`
	if since != "" && desc {
		query += " and id < $2"
	} else if since != "" && !desc {
		query += " and id > $2"
	} else if since != "" {
		query += " and id > $2"
	}
	if desc {
		query += " ORDER BY created desc, id desc"
	} else if !desc {
		query += " ORDER BY created asc, id"
	} else {
		query += " ORDER BY created, id"
	}
	return query
}

func FormQueryTreeSort(limit, threadID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE thread = $1`
			
	if since == "" {
		if desc {
			query += ` ORDER BY path DESC, id DESC`
		} else {
			query += ` ORDER BY path ASC, id ASC`
		}
	} else {
		if desc {
			query += ` AND PATH < (SELECT path FROM posts WHERE id = $2) ORDER BY path DESC, id DESC`
		} else {
			query += ` AND PATH > (SELECT path FROM posts WHERE id = $2) ORDER BY path ASC, id ASC`
		}
	}

	return query
}

func FormQueryParentTreeSort(limit, threadID int, sort, since string, desc bool) string {
	query := `SELECT id, parent, author, message, is_edited, forum, thread, created from posts
			WHERE `
			
	if since == "" {
		if desc {
			query += ` path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id DESC LIMIT $2)
			ORDER BY path[1] DESC, path, id`
		} else {
			query += ` path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL ORDER BY id ASC LIMIT $2)
			ORDER BY path ASC, id ASC`
		}
	} else {
		if desc {
			query += ` path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND PATH[1] <
				(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY path[1] DESC, path, id`
		} else {
			query += ` path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent IS NULL AND PATH[1] >
				(SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3)
			ORDER BY path ASC, id ASC`
		}
	}

	return query
}

func (tr *ThreadsRepository) ViewPosts(sort, desc, since string, limit, id int) ([]*models.Post, error) {
	postID, _ := strconv.Atoi(since)
	var desc_b bool
	if desc == "true" {
		desc_b = true
	} else {
		desc_b = false
	}
	var query  = ""
	if sort == "flat" || sort == "" {
		query = FormQueryFlatSort(limit, id, sort, since, desc_b)
		query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)
	}
	if sort == "tree" {
		query = FormQueryTreeSort(limit, id, sort, since, desc_b)
		query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", limit)
	}
	if sort == "parent_tree" {
		query = FormQueryParentTreeSort(limit, id, sort, since, desc_b)
	}
	var rows *sql.Rows
	var err error

	if since != "" {
		if sort == "parent_tree" {
			rows, err = tr.dbConn.Query(query, id, postID, limit)
		} else {
			rows, err = tr.dbConn.Query(query, id, postID)
		}
	} else {
		if sort == "parent_tree" {
			rows, err = tr.dbConn.Query(query, id, limit)
		} else {
			rows, err = tr.dbConn.Query(query, id)
		}
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	posts := make([]*models.Post, 0)
	var parent sql.NullInt64
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(&post.ID, &parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, &post.Created)

		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (tr *ThreadsRepository) UpdateVote(vote *models.Vote) (error) {
	res, _ := tr.dbConn.Exec(`UPDATE votes set voice = $1 where nickname = $2 and thread = $3 and voice !=  $1`,vote.Votes, vote.Nickname, vote.Thread)

	count, _ := res.RowsAffected()
	if count != 1 {
		return nil
	} 
	return errors.New("")
}

func (tr *ThreadsRepository) Vote(vote *models.Vote) (error) {
	_, err := tr.dbConn.Exec(`INSERT INTO votes (nickname, thread, voice)
	VALUES ($1, $2, $3)`, vote.Nickname, vote.Thread, vote.Votes)

	if err == nil {
		return nil
	} 
	return err
}

func (tr *ThreadsRepository) ThreadById(id int) (*models.Thread, error) {
	var thread models.Thread
	err := tr.dbConn.QueryRow(`SELECT slug, id, forum, title, author, message, votes, created from thread where id = $1`, id).Scan(&thread.Slug, &thread.Id, &thread.Forum, &thread.Title, &thread.Author, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return nil, err
	}
	return &thread, nil

}


func (tr *ThreadsRepository) ThreadBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread
	err := tr.dbConn.QueryRow(`SELECT slug, id, forum, title, author, message, votes, created from thread where LOWER(slug) = LOWER($1)`, slug).Scan(&thread.Slug, &thread.Id, &thread.Forum, &thread.Title, &thread.Author, &thread.Message, &thread.Votes, &thread.Created)
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

    tx, err := tr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
        return nil, err
    }

	res, err2 := tx.Query(query)
	if err2 != nil {
        rollbackErr := tx.Rollback()
        if rollbackErr != nil {
            return nil, rollbackErr
        }
        return nil, err2
    }

	newPosts := make([]*models.Post, 0)
	var parent sql.NullInt64
	defer res.Close()
	for res.Next() {
		post := &models.Post{}
		err = res.Scan(&post.ID, &parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	
		if parent.Valid {
			post.Parent = int(parent.Int64)
		}
		if err != nil {
			return nil, err
		}
		newPosts = append(newPosts, post)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

    err = tx.Commit()
    if err != nil {
        return nil, err
    }

	return newPosts, nil
}