package usecase

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"strconv"
)

type ThreadsUsecase struct {
	threadsRepo threads.ThreadsRepository
}

func NewThreadsUsecase(repo threads.ThreadsRepository) threads.ThreadsUsecase {
	return &ThreadsUsecase{
		threadsRepo: repo,
	}
}


func (tu *ThreadsUsecase) CreatePost(posts []*models.Post, slug string) ([]*models.Post, *errors.Error) {
	if len(posts) == 0 {
		return posts, nil
	}

	threadID, err := strconv.Atoi(slug)
	thread := &models.Thread{}
	if err != nil {
		thread, err = tu.threadsRepo.ThreadBySlug(slug)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
	} else {
		thread, err = tu.threadsRepo.ThreadById(threadID)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
	}

	for _, post := range posts {
		post.Thread = int(thread.Id)
		post.Forum = thread.Forum
	}

	posts, err = tu.threadsRepo.CreatePost(posts)
	return posts, nil
}