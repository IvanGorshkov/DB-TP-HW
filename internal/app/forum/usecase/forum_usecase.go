package usecase

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/forum"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"github.com/jackc/pgx"
)

type FourmUsecase struct {
	forumRepo forum.ForumRepository
	thredRepo threads.ThreadsRepository
}

func NewUserUsecase(repo forum.ForumRepository, 
	thredRepo threads.ThreadsRepository) forum.ForumUsecase {
	return &FourmUsecase{
		forumRepo: repo,
		thredRepo: thredRepo,
	}
}


func(fu *FourmUsecase) GetUserByParams(forumSlug, since, desc string, limit int) ([]*models.User, *errors.Error) {

	res, _ := fu.forumRepo.GetUserByParams(forumSlug, since, desc, limit)
	if len(res) == 0 {
		_, err := fu.forumRepo.Detail(forumSlug)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, errors.NotFoundBody("Can't find form with slug " + forumSlug + "\n")
			}
			return nil, errors.UnexpectedInternal(err)
		}
	}

	return res, nil
}

func(fu *FourmUsecase) GetThreadsByParams(forumSlug, since, desc string, limit int) ([]*models.Thread, *errors.Error) {
	res, _ := fu.forumRepo.GetThreadsByParams(forumSlug, since, desc, limit)
	if len(res) == 0 {
		_, err := fu.forumRepo.Detail(forumSlug)
		if err != nil {
			return nil, errors.NotFoundBody("Can't find user with nickname " + forumSlug + "\n")
		}

		return res,nil
	}

	return res, nil
}

func(fu *FourmUsecase) CreateThread(thread *models.Thread) (*models.Thread, *errors.Error) {
	res, err := fu.forumRepo.ThreadCreate(thread) 
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23503" {
			return res, errors.NotFoundBody("Can't find user with nickname " + thread.Author + "\n")
		}

		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23502" {
			return res, errors.NotFoundBody("Can't find thread forum by slug " + thread.Forum + "\n")
		}

		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == "23505" {
			thr, _ := fu.thredRepo.ThreadBySlug(thread.Slug)

			return thr, errors.CustomErrors[errors.ConflictError]
		}

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
		if err == pgx.ErrNoRows {
			return nil, errors.NotFoundBody("Can't find form with slug " + slug + "\n")
		}
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}
