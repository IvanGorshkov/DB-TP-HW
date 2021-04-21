package forum

import "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"


type ForumRepository interface {
	Create(forum *models.Forum) (*models.Forum, error)
	Detail(slug string) (*models.Forum, error)
	ThreadCreate(thread *models.Thread) (*models.Thread, error)
	GetThreadsByParams(forumSlug, since, desc string, limit int) ([]*models.Thread, error)
}
