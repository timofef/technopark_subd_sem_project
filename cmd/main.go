package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/jackc/pgx"
	handlers "github.com/timofef/technopark_subd_sem_project/delivery"
	repository "github.com/timofef/technopark_subd_sem_project/repository/implementations"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/implementations"
	"github.com/valyala/fasthttp"
)

func main() {
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			User:     "admin",
			Password: "postgres",
			Port:     5432,
			Database: "forum",
		},
		MaxConnections: 50,
	})
	if err != nil {
		fmt.Println(err)
	}

	router := fasthttprouter.New()

	userRepo := repository.NewUserRepo(pool)
	forumRepo := repository.NewForumRepo(pool)
	threadRepo := repository.NewThreadRepo(pool)
	postRepo := repository.NewPostRepo(pool)

	userUsecase := usecase.NewUserUsecase(userRepo)
	forumUsecase := usecase.NewForumUsecase(forumRepo, userRepo, threadRepo)
	threadUsecase := usecase.NewThreadUsecase(threadRepo, forumRepo, postRepo, userRepo)

	handlers.NewUserHandler(router, userUsecase)
	handlers.NewForumHandler(router, forumUsecase)
	handlers.NewThreadHandler(router, threadUsecase)

	fmt.Println("http server started on 5000 port")

	err = fasthttp.ListenAndServe(":5000", router.Handler)
}
