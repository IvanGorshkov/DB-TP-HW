package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/service"
	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	serviceUsecase service.ServiceUsecase
}

func NewServiceHandler(serviceUsecase service.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUsecase: serviceUsecase,
	}
}


func (sh *ServiceHandler) Configure(r *mux.Router) {
	r.HandleFunc("/service/status", sh.getStatus).Methods(http.MethodGet)
}

func (sh *ServiceHandler) getStatus(w http.ResponseWriter, r *http.Request) {
	res, err := sh.serviceUsecase.getStatus()

	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jData, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(jData)
}
