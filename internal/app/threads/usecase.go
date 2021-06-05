package threads

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type ThreadsUsecase interface {
	Update(thread *models.Thread, slug string) (*models.Thread, *errors.Error)
	CreatePost(posts []models.Post, slug string) ([]models.Post, *errors.Error)
	VoteByIdOrSlag(vote *models.Vote, slug string) (*models.Thread, *errors.Error)
	Detail(slug_or_id string) (*models.Thread, *errors.Error)
	ViewPosts(id, sort, desc, since string, limit int)  ([]models.Post, *errors.Error)
} 
