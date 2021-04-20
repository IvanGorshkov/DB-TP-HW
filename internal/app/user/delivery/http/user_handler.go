package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUsecase user.UserUsecase
}

func NewUserHandler(userUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}


func (uh *UserHandler) Configure(r *mux.Router) {
	r.HandleFunc("/user/{nickname}/create", uh.createUser).Methods(http.MethodPost)
}

func (uh *UserHandler) createUser(w http.ResponseWriter, r *http.Request){


	userModel := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	userModel.Nickname = vars["nickname"]
	usr, err := uh.userUsecase.Create(userModel)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		fmt.Println(err)
		return
	}
}