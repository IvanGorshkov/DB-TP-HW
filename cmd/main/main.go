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

	threadHandler "github.com/IvanGorshkov/DB-TP-HW/internal/app/threads/delivery/http"
	threadRepo "github.com/IvanGorshkov/DB-TP-HW/internal/app/threads/repository/postgres"
	threadUsecase "github.com/IvanGorshkov/DB-TP-HW/internal/app/threads/usecase"

	postHandler "github.com/IvanGorshkov/DB-TP-HW/internal/app/post/delivery/http"
	postRepo "github.com/IvanGorshkov/DB-TP-HW/internal/app/post/repository/postgres"
	postUsecase "github.com/IvanGorshkov/DB-TP-HW/internal/app/post/usecase"

	serviceHandler "github.com/IvanGorshkov/DB-TP-HW/internal/app/service/delivery/http"
	serviceRepo "github.com/IvanGorshkov/DB-TP-HW/internal/app/service/repository/postgres"
	serviceUsecase "github.com/IvanGorshkov/DB-TP-HW/internal/app/service/usecase"
)


func main() {
	postgresDB, err := databases.NewPostgres(databases.GetPostgresConfig())
	if err != nil {

		return
	}

	defer postgresDB.Close()

	router := mux.NewRouter()

	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	userUsecase := userUsecase.NewUserUsecase(userRepo)
	userHandler := userHandler.NewUserHandler(userUsecase)


	threadRepo := threadRepo.NewThreadsRepository(postgresDB.GetDatabase())
	threadUsecase := threadUsecase.NewThreadsUsecase(threadRepo)
	threadHandler := threadHandler.NewThreadsHandler(threadUsecase)

	forumRepo := forumRepo.NewForumRepository(postgresDB.GetDatabase())
	forumUsecase := forumUsecase.NewUserUsecase(forumRepo, threadRepo)
	forumHandler := forumHandler.NewForumHandler(forumUsecase)


	
	postRepo := postRepo.NewPostRepository(postgresDB.GetDatabase())
	postUsecase := postUsecase.NewThreadsUsecase(postRepo, userRepo, forumRepo, threadRepo)
	postHandler := postHandler.NewThreadsHandler(postUsecase)


	serviceRepo := serviceRepo.NewServiceRepository(postgresDB.GetDatabase())
	serviceUsecase := serviceUsecase.NewUserUsecase(serviceRepo)
	serviceHandler := serviceHandler.NewServiceHandler(serviceUsecase)

	api := router.PathPrefix("/api/").Subrouter()

	userHandler.Configure(api)
	forumHandler.Configure(api)
	threadHandler.Configure(api)
	postHandler.Configure(api)
	serviceHandler.Configure(api)
	
	server := http.Server{
		Addr:         fmt.Sprint(":5000"),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	_ = server.ListenAndServe()
}