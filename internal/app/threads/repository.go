package threads

import "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"


type ThreadsRepository interface {
	ThreadById(id int) (*models.Thread, error)
	ThreadBySlug(slug string) (*models.Thread, error)
	CreatePost(posts []*models.Post) ([]*models.Post, error) 
}
