package threads

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type ThreadsUsecase interface {
	CreatePost(posts []*models.Post, slug string) ([]*models.Post, *errors.Error)
} 
