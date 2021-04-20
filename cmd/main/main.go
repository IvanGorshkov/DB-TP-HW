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


	forumHandler "github.com/IvanGorshkov/DB-TP-HW/internal/app/forum/delivery/http"
	forumRepo "github.com/IvanGorshkov/DB-TP-HW/internal/app/forum/repository/postgres"
	forumUsecase "github.com/IvanGorshkov/DB-TP-HW/internal/app/forum/usecase"
)


func main() {
	postgresDB, err := databases.NewPostgres(databases.GetPostgresConfig())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer postgresDB.Close()

	router := mux.NewRouter()

	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	userUsecase := userUsecase.NewUserUsecase(userRepo)
	userHandler := userHandler.NewUserHandler(userUsecase)

	forumRepo := forumRepo.NewForumRepository(postgresDB.GetDatabase())
	forumUsecase := forumUsecase.NewUserUsecase(forumRepo)
	forumHandler := forumHandler.NewForumHandler(forumUsecase)

	api := router.PathPrefix("/api/").Subrouter()

	userHandler.Configure(api)
	forumHandler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":5000"),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	_ = server.ListenAndServe()
}