package main

import (
	"github.com/jackc/pgx"
	"github.com/buaazp/fasthttprouter"
)

func main() {
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			User:     "docker",
			Password: "docker",
			Port:     5432,
			Database: "docker",
		},
		MaxConnections: 50,
	})

	router := fasthttprouter.New()



}
