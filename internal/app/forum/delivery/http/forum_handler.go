package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/forum"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/gorilla/mux"
)

type ForumHandler struct {
	forumUsecase forum.ForumUsecase
}

func NewForumHandler(forumUsecase forum.ForumUsecase) *ForumHandler {
	return &ForumHandler{
		forumUsecase: forumUsecase,
	}
}


func (fh *ForumHandler) Configure(r *mux.Router) {
	r.HandleFunc("/forum/create", fh.createForum).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/details", fh.detailsForum).Methods(http.MethodGet)
}

func (fh *ForumHandler) detailsForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err2 := fh.forumUsecase.Detail(vars["slug"])
	if  err2 != nil {
		if  err2.ErrorCode == errors.InternalError {
			return
		}

		if err2.ErrorCode == errors.NotFoundError {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err2.HttpError)
			messagee := errors.Message{ Message: err2.Message}
			err := json.NewEncoder(w).Encode(messagee)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (fh *ForumHandler) createForum(w http.ResponseWriter, r *http.Request) {
	forumModel := &models.Forum{}
	err := json.NewDecoder(r.Body).Decode(&forumModel)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err2 := fh.forumUsecase.Create(forumModel)
	if  err2 != nil {
		if  err2.ErrorCode == errors.InternalError {
			return
		}

		if err2.ErrorCode == errors.ConflictError {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err2.HttpError)
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		if err2.ErrorCode == errors.NotFoundError {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err2.HttpError)
			messagee := errors.Message{ Message: err2.Message}
			err := json.NewEncoder(w).Encode(messagee)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}
}