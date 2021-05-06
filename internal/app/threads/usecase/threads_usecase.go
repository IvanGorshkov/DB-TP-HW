package usecase

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"github.com/jackc/pgx"
)

type ThreadsUsecase struct {
	threadsRepo threads.ThreadsRepository
}

func NewThreadsUsecase(repo threads.ThreadsRepository) threads.ThreadsUsecase {
	return &ThreadsUsecase{
		threadsRepo: repo,
	}
}


func (tu *ThreadsUsecase) Update(thread *models.Thread, slug string) (*models.Thread, *errors.Error) {
	threadID, err := strconv.Atoi(slug)
	if err != nil {
		_, err = tu.threadsRepo.ThreadBySlug(slug)

		if err != nil {
				if err == sql.ErrNoRows  {
					return nil, errors.NotFoundBody("Can't find thread by slug: " + slug + "\n" )
				}
			return nil, errors.UnexpectedInternal(err)
		}

		thread, err = tu.threadsRepo.UpdateBySlug(thread,slug)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}

	} else {
		_, err = tu.threadsRepo.ThreadById(threadID)

		if err != nil {
				if err == sql.ErrNoRows  {
					return nil, errors.NotFoundBody("Can't find thread by slug: " + slug + "\n" )
				}
			return nil, errors.UnexpectedInternal(err)
		}


		thread, err = tu.threadsRepo.UpdateById(thread, threadID)
		if err != nil {
			return nil, errors.UnexpectedInternal(err)
		}
	}

	return thread, nil
}

func (tu *ThreadsUsecase) ViewPosts(id, sort, desc, since string, limit int) ([]*models.Post, *errors.Error) {
	threadID, err := strconv.Atoi(id)
	var thread = &models.Thread{}
	if err != nil {
		thread, err = tu.threadsRepo.ThreadBySlug(id)

		if err != nil {
				if err == sql.ErrNoRows  {
					return nil, errors.NotFoundBody("Can't find thread by slug: " + id + "\n" )
				}
			return nil, errors.UnexpectedInternal(err)
		}
	} else {
		thread, err = tu.threadsRepo.ThreadById(threadID)

		if err != nil {
			if err == sql.ErrNoRows  {
				return nil, errors.NotFoundBody("Can't find thread by slug: " + id + "\n" )
			}
			return nil, errors.UnexpectedInternal(err)
		}
		
	}

	res, err2 := tu.threadsRepo.ViewPosts(sort, desc, since, limit, int(thread.Id))
	if err2 != nil { 
		return nil, errors.UnexpectedInternal(err2)
	}
	return res, nil
}

func (tu *ThreadsUsecase) Detail(slug_or_id string) (*models.Thread, *errors.Error) {
	threadID, err := strconv.Atoi(slug_or_id)
	var thread = &models.Thread{}
	if err != nil {
		thread, err = tu.threadsRepo.ThreadBySlug(slug_or_id)

		if err != nil {
				if err == sql.ErrNoRows  {
					return nil, errors.NotFoundBody("Can't find thread by slug: " + slug_or_id + "\n" )
				}
			return nil, errors.UnexpectedInternal(err)
		}
	} else {
		thread, err = tu.threadsRepo.ThreadById(threadID)

		if err != nil {
			if err == sql.ErrNoRows  {
				return nil, errors.NotFoundBody("Can't find thread by slug: " + slug_or_id + "\n" )
			}
			return nil, errors.UnexpectedInternal(err)
		}
		
	}
	return thread, nil
}

func (tu *ThreadsUsecase) VoteByIdOrSlag(vote *models.Vote, slug string) (*models.Thread, *errors.Error) {
	threadID, err := strconv.Atoi(slug)
	var thread = &models.Thread{}
	if err != nil {
		thread, err = tu.threadsRepo.ThreadBySlug(slug)
		if err != nil {
				if err == sql.ErrNoRows  {
					return nil, errors.NotFoundBody("Can't find thread by slug: " + slug + "\n" )
				}
			return nil, errors.UnexpectedInternal(err)
		}
	} else {
		thread, err = tu.threadsRepo.ThreadById(threadID)
		if err != nil {
			if err == sql.ErrNoRows  {
				return nil, errors.NotFoundBody("Can't find thread by slug: " + slug + "\n" )
			}
			return nil, errors.UnexpectedInternal(err)
		}
	}
	vote.Thread = int(thread.Id)
	
	err = tu.threadsRepo.Vote(vote)
	if err != nil {
		if err.(pgx.PgError).Code == "23503" {
			return nil, errors.NotFoundBody(vote.Nickname)
		}
		if err.(pgx.PgError).Code == "23505" {
			err = tu.threadsRepo.UpdateVote(vote)
			if err != nil {
				if vote.Votes < 0 {
					thread.Votes += 2 * int32(vote.Votes);
				} else {
					thread.Votes += int32(vote.Votes);
				}
			} 

			return thread, nil
		}

		return nil, errors.UnexpectedInternal(err)
	}
	thread.Votes += int32(vote.Votes);
	return thread, nil
}

func (tu *ThreadsUsecase) CreatePost(posts []*models.Post, slug string) ([]*models.Post, *errors.Error) {

	threadID, err := strconv.Atoi(slug)
	var thread = &models.Thread{}
	if err != nil {
		thread, err = tu.threadsRepo.ThreadBySlug(slug)
		fmt.Println(err, ' ', thread)
		if err != nil {
			if err == sql.ErrNoRows  {
				return nil, errors.NotFoundBody("Can't find thread by slug: " + slug + "\n" )
			}
			return nil, errors.UnexpectedInternal(err)
		}
	} else {
		thread, err = tu.threadsRepo.ThreadById(threadID)
		fmt.Println(err, ' ', thread)
		if err != nil {
			if err == sql.ErrNoRows  {
				return nil, errors.NotFoundBody("Can't find post thread by id: " + slug + "\n" )
			}
			return nil, errors.UnexpectedInternal(err)
		}
	}

	if len(posts) == 0 {
		return posts, nil
	}

	
	for _, post := range posts {
		post.Thread = int(thread.Id)
		post.Forum = thread.Forum
	}

	posts, err = tu.threadsRepo.CreatePost(posts)
	fmt.Println(err)
	if err != nil {
		if err.Error() == "409" {
			return nil , errors.ConflictErrorBody("Parent post was created in another thread")
		}

		if err.Error() == "404" {
			return nil, errors.NotFoundBody(`Can't find post author by nickname: \n`)
		}
	}
	return posts, nil
}