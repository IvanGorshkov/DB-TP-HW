package post

import "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"



type PostRepository interface {
	GetPostById(id int) (*models.Post, error)
	Update(id int, post models.Post) (*models.Post, error)
}
