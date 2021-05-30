package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	r.HandleFunc("/forum/{slug}/create", fh.CreateThread).Methods(http.MethodPost)
	r.HandleFunc("/forum/{slug}/threads", fh.GetThreads).Methods(http.MethodGet)   
	r.HandleFunc("/forum/{slug}/users", fh.GetUsers).Methods(http.MethodGet)   
}

func (fh *ForumHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/forum/{slug}/users Get")
	vars := mux.Vars(r)

	limit, err := strconv.Atoi(string(r.FormValue("limit")))
	if err != nil {
		fmt.Println(err)
	}
	since := string(r.FormValue("since"))
	if err != nil {
		fmt.Println(err)
	}
	desc := string(r.FormValue("desc"))
	if err != nil {
		fmt.Println(err)
	}

	res, err2 := fh.forumUsecase.GetUserByParams(vars["slug"], since, desc, limit)

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
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (fh *ForumHandler) GetThreads(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/forum/{slug}/threads Get")
	vars := mux.Vars(r)

	limit, err := strconv.Atoi(string(r.FormValue("limit")))
	if err != nil {
		fmt.Println(err)
	}
	since := string(r.FormValue("since"))
	if err != nil {
		fmt.Println(err)
	}
	desc := string(r.FormValue("desc"))
	if err != nil {
		fmt.Println(err)
	}

	res, err2 := fh.forumUsecase.GetThreadsByParams(vars["slug"], since, desc, limit)

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
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (fh *ForumHandler) CreateThread(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/forum/{slug}/create POST")
	threadModel := &models.Thread{}
	err := json.NewDecoder(r.Body).Decode(&threadModel)
	if err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	threadModel.Forum = vars["slug"]
	res, err2 := fh.forumUsecase.CreateThread(threadModel)

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

		if err2.ErrorCode == errors.ConflictError {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			err = json.NewEncoder(w).Encode(res)
			if err != nil {
				fmt.Println(err)
				return
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

func (fh *ForumHandler) detailsForum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/forum/{slug}/details GET")
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
	fmt.Println("/forum/create POST")

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