package threads

import "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"


type ThreadsRepository interface {
	ThreadById(id int) (*models.Thread, error)
	ThreadBySlug(slug string) (*models.Thread, error)
	UpdateById(thread *models.Thread, id int) (*models.Thread, error)
	UpdateBySlug(thread *models.Thread, slug string) (*models.Thread, error)
	CreatePost(posts []*models.Post) ([]*models.Post, error) 
	Vote(vote *models.Vote) (error)
	UpdateVote(vote *models.Vote) (error)
	ViewPosts(sort, desc, since string, limit, id int) ([]*models.Post, error)
}
