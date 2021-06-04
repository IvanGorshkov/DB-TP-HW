package repository

import (
    "errors"
    "fmt"
    "github.com/go-openapi/strfmt"
    "github.com/jackc/pgx"
    "time"

    "github.com/IvanGorshkov/DB-TP-HW/internal/app/forum"
    "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type ForumRepository struct {
    dbConn *pgx.ConnPool
}

func NewForumRepository(conn *pgx.ConnPool) forum.ForumRepository {
    return &ForumRepository{
        dbConn: conn,
    }
}

func(fr *ForumRepository) GetUserByParams(forumSlug, since, desc string, limit int) ([]*models.User, error) {
	query := `SELECT u.nickname, u.email, u.fullname, u.about from users_forum
    join users u on users_forum.nickname = u.nickname
    where slug = $1 `

	if since != "" && desc == "true" {
		query += fmt.Sprintf(` and CAST(LOWER(u.nickname) AS bytea) < CAST(LOWER('%s') AS bytea)`, since)
	} else if since != "" {
		query += fmt.Sprintf(` and CAST(LOWER(u.nickname) AS bytea) > CAST(LOWER('%s') AS bytea)`, since)
	}

    if desc == "true" {
        query += ` order BY CAST(LOWER(u.nickname) AS bytea) desc`
    } else if desc == "false" {
        query += ` order BY CAST(LOWER(u.nickname) AS bytea) asc`
    } else {
        query += ` order BY CAST(LOWER(u.nickname) AS bytea) asc`
    }

	query += " LIMIT NULLIF($2, 0)"

    q, err := fr.dbConn.Query(query, forumSlug, limit)
    defer  q.Close()
    if err != nil {

        return nil, err
    }
	users := make([]*models.User, 0)
    for q.Next() {
        user := &models.User{}
        err = q.Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About)
            if err != nil {
                return nil, err
            }
			users = append(users, user)
    }


    return users, nil
}

func(fr *ForumRepository) GetThreadsByParams(forumSlug, since, desc string, limit int) ([]*models.Thread, error) {
    query := `SELECT t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from thread as t
    LEFT JOIN forum f on t.forum = f.slug where f.slug = $1`

    if since != "" && desc == "true" {
        query += " and t.created <= $2"
    } else if since != "" && desc == "false" {
        query += " and t.created >= $2"
    } else if since != "" {
        query += " and t.created >= $2"
    }

    if desc == "true" {
        query += " ORDER BY t.created desc"
    } else if desc == "false" {
        query += " ORDER BY t.created asc"
    } else {
        query += " ORDER BY t.created"
    }

    var args []interface{}
    if since != "" {
        query += " LIMIT NULLIF($3, 0)"
        args = append(args, forumSlug, since, limit)
    } else {
        query += " LIMIT NULLIF($2, 0)"
        args = append(args, forumSlug, limit)
    }

    q, err := fr.dbConn.Query(query, args...)
    defer q.Close()
    if err != nil {

        return nil, err
    }

    threads := make([]*models.Thread, 0)
    for q.Next() {
        thread := &models.Thread{}
        var created time.Time
        err = q.Scan(&thread.Id, &thread.Title,
            &thread.Author, &thread.Forum, &thread.Message, &thread.Votes,
            &thread.Slug, &created)
            if err != nil {
                return nil, err
            }

        thread.Created = strfmt.DateTime(created.UTC()).String()
        threads = append(threads, thread)
    }

    return threads, nil
}

func(fr *ForumRepository) ThreadCreate(thread *models.Thread) (*models.Thread, error) {

    tx, err := fr.dbConn.Begin()
    if err != nil {
        return nil, err
    }

    if thread.Created == "" {
        err = tx.QueryRow(`
            INSERT INTO thread (title, author, forum, message, slug) VALUES ($1, $2, (SELECT slug from forum where slug = $3), $4, $5) returning id, forum
        `, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug).Scan(&thread.Id, &thread.Forum) 
    } else {
        err = tx.QueryRow(`
            INSERT INTO thread (title, author, forum, message, created, slug) VALUES ($1, $2, (SELECT slug from forum where slug = $3), $4, $5, $6) returning id, forum
        `, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Created, thread.Slug).Scan(&thread.Id, &thread.Forum) 
    }

    if err != nil {
        rollbackErr := tx.Rollback()
        if rollbackErr != nil {
            return nil, rollbackErr
        }
        return nil, err
    }

    err = tx.Commit()
    if err != nil {
        return nil, err
    }
    return thread, nil
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
    
    queryForum, err := fr.dbConn.Query(`SELECT title, nickname, slug, post, threads from forum where slug = $1`,forum.Slug)
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
    err = fr.dbConn.QueryRow(`SELECT nickname from users where nickname = $1`,forum.User).Scan(&user)
    forum.User = user
    tx, err := fr.dbConn.Begin()
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
        return nil, err
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
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, pgx.ErrNoRows
        }
        return nil, err
    }

    return &forum, nil
}
