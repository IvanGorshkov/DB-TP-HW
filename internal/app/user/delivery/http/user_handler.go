package delivery

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/user"
	"github.com/gorilla/mux"
	"net/http"
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
	r.HandleFunc("/user/{nickname}/create", uh.createUser).Methods(http.MethodGet)
}

func (uh *UserHandler) createUser(w http.ResponseWriter, r *http.Request){
	
}