package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
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
	r.HandleFunc("/user/{nickname}/profile", uh.getUser).Methods(http.MethodGet)
	r.HandleFunc("/user/{nickname}/profile", uh.updateUser).Methods(http.MethodPost)
}

func (uh *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usr, err2 := uh.userUsecase.GetProfile(vars["nickname"])
	if  err2 != nil {
		if  err2.ErrorCode == errors.InternalError {
			return
		}
		if err2.ErrorCode == errors.NotFoundError {
			w.WriteHeader(err2.HttpError)
			w.Header().Set("Content-Type", "application/json")
			messagee := errors.Message{ Message: err2.Message}
			err := json.NewEncoder(w).Encode(messagee)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(usr)
	if err != nil {
		fmt.Println(err)
		return
	}
}


func (uh *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userModel := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	userModel.Nickname = vars["nickname"]
	usr, err2 := uh.userUsecase.UpdateProfile(userModel)
	if  err2 != nil {
		if  err2.ErrorCode == errors.InternalError {
			return
		}

		if err2.ErrorCode == errors.NotFoundError ||  err2.ErrorCode == errors.ConflictError {
			w.WriteHeader(err2.HttpError)
			w.Header().Set("Content-Type", "application/json")
			messagee := errors.Message{ Message: err2.Message}
			err := json.NewEncoder(w).Encode(messagee)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (uh *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	userModel := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		fmt.Println(err)
		return
	}
	vars := mux.Vars(r)
	userModel.Nickname = vars["nickname"]
	usr, err2 := uh.userUsecase.Create(userModel)
	if  err2 != nil {
		if  err2.ErrorCode == errors.InternalError {
			return
		}

		if err2.ErrorCode == errors.ConflictError {
			w.WriteHeader(err2.HttpError)
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(usr)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(usr[0])
	if err != nil {
		fmt.Println(err)
		return
	}
}