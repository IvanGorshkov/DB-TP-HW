package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"github.com/gorilla/mux"
)

type ThreadsHandler struct {
	threadsUsecase threads.ThreadsUsecase
}

func NewThreadsHandler(threadsUsecase threads.ThreadsUsecase) *ThreadsHandler {
	return &ThreadsHandler{
		threadsUsecase: threadsUsecase,
	}
}


func (th *ThreadsHandler) Configure(r *mux.Router) {
	r.HandleFunc("/thread/{slug}/create", th.postsCreate).Methods(http.MethodPost)
	r.HandleFunc("/thread/{slug_or_id}/vote", th.Vote).Methods(http.MethodPost)
	r.HandleFunc("/thread/{slug_or_id}/details", th.Detail).Methods(http.MethodGet)
	r.HandleFunc("/thread/{slug_or_id}/posts", th.ViewPosts).Methods(http.MethodGet)
	r.HandleFunc("/thread/{slug_or_id}/details", th.Update).Methods(http.MethodPost)
}


func (th *ThreadsHandler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/thread/{slug_or_id}/details Post")
	vars := mux.Vars(r)
	var thread = models.Thread{}
	err := json.NewDecoder(r.Body).Decode(&thread)
	res, err2 := th.threadsUsecase.Update(&thread, vars["slug_or_id"])
	if err2 != nil {
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (th *ThreadsHandler) ViewPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/thread/{slug_or_id}/posts Get")

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

	sort := string(r.FormValue("sort"))
	if err != nil {
		fmt.Println(err)
	}

	res, err2 := th.threadsUsecase.ViewPosts(vars["slug_or_id"], sort, desc, since, limit)

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

func (th *ThreadsHandler) Detail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/thread/{slug_or_id}/details Get")
	vars := mux.Vars(r)
	res, err2 := th.threadsUsecase.Detail(vars["slug_or_id"])
	if err2 != nil {
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (th *ThreadsHandler) Vote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/thread/{slug_or_id}/vote Post")
	var vote = models.Vote{}
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		return
	}
	vars := mux.Vars(r)

	res, err2 := th.threadsUsecase.VoteByIdOrSlag(&vote, vars["slug_or_id"])
	if err2 != nil {
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (th *ThreadsHandler) postsCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/thread/{slug}/create Post")
	posts := make([]*models.Post, 0)
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		return
	}
	vars := mux.Vars(r)
	res, err2 := th.threadsUsecase.CreatePost(posts, vars["slug"])
	if err2 != nil {
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