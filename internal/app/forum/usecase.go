package forum

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type ForumUsecase interface {
	Create(forum *models.Forum) (*models.Forum, *errors.Error)
	Detail(slug string) (*models.Forum, *errors.Error)
	CreateThread(thread *models.Thread) (*models.Thread, *errors.Error)
	GetThreadsByParams(forumSlug, since, desc string, limit int) ([]*models.Thread, *errors.Error)
	GetUserByParams(forumSlug, since, desc string, limit int) ([]*models.User, *errors.Error)
}
