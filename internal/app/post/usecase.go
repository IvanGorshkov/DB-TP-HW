package post

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)


type PostUsecase interface {
	Detail(id int, related []string) (*models.PostFull, *errors.Error)
	Update(id int, post models.Post) (*models.Post, *errors.Error)
} 
