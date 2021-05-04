package usecase

import (
	"database/sql"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/forum"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)

type FourmUsecase struct {
	forumRepo forum.ForumRepository
}

func NewUserUsecase(repo forum.ForumRepository) forum.ForumUsecase {
	return &FourmUsecase{
		forumRepo: repo,
	}
}


func(fu *FourmUsecase) GetUserByParams(forumSlug, since, desc string, limit int) ([]*models.User, *errors.Error) {
	_, err := fu.forumRepo.Detail(forumSlug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundBody("Can't find form with slug " + forumSlug + "\n")
		}
		return nil, errors.UnexpectedInternal(err)
	}
	res, err2 := fu.forumRepo.GetUserByParams(forumSlug, since, desc, limit)
	if err2 != nil {
		return nil, errors.UnexpectedInternal(err2)
	}

	return res, nil
}

func(fu *FourmUsecase) GetThreadsByParams(forumSlug, since, desc string, limit int) ([]*models.Thread, *errors.Error) {
	_, err := fu.forumRepo.Detail(forumSlug)
	if err != nil {
		return nil, errors.NotFoundBody("Can't find user with nickname " + forumSlug + "\n")
	}

	res, err := fu.forumRepo.GetThreadsByParams(forumSlug, since, desc, limit) 
	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}

func(fu *FourmUsecase) CreateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	res, err := fu.forumRepo.ThreadCreate(thread) 
	if err != nil {
		if err.Error() == "409" {
			return res, errors.CustomErrors[errors.ConflictError]
		}
		if err.Error() == "404" {
			return res, errors.NotFoundBody("Can't find user with nickname " + thread.Author + "\n")
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}

func (fu *FourmUsecase) Create(forum *models.Forum) (*models.Forum, *errors.Error) {
	res, err := fu.forumRepo.Create(forum) 
	if err != nil {
		if err.Error() == "409" {
			return res, errors.CustomErrors[errors.ConflictError]
		}
		if err.Error() == "404" {
			return res, errors.NotFoundBody("Can't find user with nickname " + forum.User + "\n")
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}

func (fu *FourmUsecase) Detail(slug string) (*models.Forum, *errors.Error) {
	res, err := fu.forumRepo.Detail(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundBody("Can't find form with slug " + slug + "\n")
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}
