package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"time"

	"github.com/IvanGorshkov/DB-TP-HW/internal/app/databases"
	userHandler "github.com/IvanGorshkov/DB-TP-HW/internal/app/user/delivery/http"
	userRepo "github.com/IvanGorshkov/DB-TP-HW/internal/app/user/repository/postgres"
	userUsecase "github.com/IvanGorshkov/DB-TP-HW/internal/app/user/usecase"
)


func main() {
	postgresDB, err := databases.NewPostgres(databases.GetPostgresConfig())
	if err != nil {
		return
	}

	defer postgresDB.Close()

	router := mux.NewRouter()
	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	userUsecase := userUsecase.NewProductUsecase(userRepo)
	userHandler := userHandler.NewUserHandler(userUsecase)
	api := router.PathPrefix("/api/").Subrouter()
	userHandler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":5000"),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	_ = server.ListenAndServe()
}