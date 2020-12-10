package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/buaazp/fasthttprouter"
	repository "github.com/timofef/technopark_subd_sem_project/repository/implementations"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/implementations"
	handlers "github.com/timofef/technopark_subd_sem_project/delivery"
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
	if err = userRepo.PrepareStatements(); err != nil {
		fmt.Println(err)
	}
	userUsecase := usecase.NewUserUsecase(userRepo)
	_ = handlers.NewUserHandler(router, userUsecase)

	fmt.Println("http server started on 5000 port: ")
	err = fasthttp.ListenAndServe(":5000", router.Handler)
}
