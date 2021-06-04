package usecase

import (
	"github.com/jackc/pgx"
	"strconv"
	"strings"

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
		if err == pgx.ErrNoRows  {
			return nil, errors.NotFoundBody("Can't find thread by slug: " + strconv.Itoa(id) + "\n" )
		}

		return nil, errors.UnexpectedInternal(err)
	}
	return res, nil
}

func (pu *PostUsecase) Detail(id int, related []string) (*models.PostFull, *errors.Error) {
	var postFull, _ =  pu.postRepository.GetPostFullbyId(id, strings.Join(related, ""))

	if postFull == nil {
		_, err := pu.postRepository.GetPostById(id)
		if err == pgx.ErrNoRows  {
			return nil, errors.NotFoundBody("Can't find thread by slug: " + strconv.Itoa(id) + "\n" )
		}
	}


	return postFull, nil
}