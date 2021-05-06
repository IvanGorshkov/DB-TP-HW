package repository

import (
	"database/sql"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/service"
)

type ServiceRepository struct {
	dbConn *sql.DB
}

func NewServiceRepository(conn *sql.DB) service.ServiceRepository {
	return &ServiceRepository{
		dbConn: conn,
	}
}

 
func(sr *ServiceRepository) GetStatus() (*models.Status, error) {
	var status models.Status
	err := sr.dbConn.QueryRow(
		`SELECT * FROM 
		(SELECT COUNT(*) FROM forum) as forumCount,
		(SELECT COUNT(*) FROM posts) as postCount,
		(SELECT COUNT(*) FROM thread) as threadCount,
		(SELECT COUNT(*) FROM users) as usersCount;`,
	).Scan(&status.Forum, &status.Post, &status.Thread, &status.User)

	if err != nil {
		return nil, err
	}
	return &status, nil
}