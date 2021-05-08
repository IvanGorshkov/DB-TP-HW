package usecase

import (
	"database/sql"
	"strconv"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/forum"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/post"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
)

type PostUsecase struct {
	postRepository post.PostRepository
	userRepository user.UserRepository
	forumRepository forum.ForumRepository
	threadsRepository threads.ThreadsRepository
}

func NewThreadsUsecase(
	postRepository post.PostRepository,
	userRepository user.UserRepository,
	forumRepository forum.ForumRepository,
	threadsRepository threads.ThreadsRepository) post.PostUsecase{
	return &PostUsecase{
		postRepository: postRepository,
		userRepository: userRepository,
		forumRepository: forumRepository,
		threadsRepository: threadsRepository,
	}
}


func (pu *PostUsecase)  Update(id int, post models.Post) (*models.Post, *errors.Error) {
	res, err := pu.postRepository.Update(id, post)
	if err != nil {
		if err == sql.ErrNoRows  {
			return nil, errors.NotFoundBody("Can't find thread by slug: " + strconv.Itoa(id) + "\n" )
		}

		return nil, errors.UnexpectedInternal(err)
	}
	return res, nil
}

func (pu *PostUsecase) Detail(id int, related []string) (*models.PostFull, *errors.Error) {
	var postFull models.PostFull
	post, err := pu.postRepository.GetPostById(id)
	if err != nil {
		if err == sql.ErrNoRows  {
			return nil, errors.NotFoundBody("Can't find thread by slug: " + strconv.Itoa(id) + "\n" )
		}

		return nil, errors.UnexpectedInternal(err)
	}
	
	postFull.Post = *post
	
	for _, item := range related {  
		switch item {
			case "user": {
				user, err := pu.userRepository.GetProfile(post.Author)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}
				postFull.Author = user
			}
			case "forum": {
				forum, err := pu.forumRepository.Detail(post.Forum)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}
				postFull.Forum = forum
			}
			case "thread": {
				thread, err := pu.threadsRepository.ThreadById(post.Thread)
				if err != nil {
					return nil, errors.UnexpectedInternal(err)
				}
				postFull.Thread = thread
			}
		}
	}

	return &postFull, nil
}