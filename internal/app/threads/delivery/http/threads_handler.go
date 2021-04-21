package delivery

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/threads"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"encoding/json"
	"fmt"
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
}

func (th *ThreadsHandler) postsCreate(w http.ResponseWriter, r *http.Request) {
	posts := make([]*models.Post, 0)
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		return
	}
	vars := mux.Vars(r)
	res, err2 := th.threadsUsecase.CreatePost(posts, vars["slug"])
	if err2 != nil {

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println(err)
		return
	}

}