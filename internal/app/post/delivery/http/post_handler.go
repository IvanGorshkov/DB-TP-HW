package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/post"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	postUsecase post.PostUsecase
}

func NewThreadsHandler(postUsecase post.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

func (ph *PostHandler) Configure(r *mux.Router) {
	r.HandleFunc("/post/{id}/details", ph.Detail).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/details", ph.Update).Methods(http.MethodPost)
}

func (ph *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err)
	}

	var post = models.Post{}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		return
	}

	res, err2 := ph.postUsecase.Update(postID, post)
	if err2 != nil {
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

func (ph *PostHandler) Detail(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err)
	}
	array := strings.Split(r.FormValue("related"), ",") 

	res, err2 := ph.postUsecase.Detail(postID, array)
	if err2 != nil {
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