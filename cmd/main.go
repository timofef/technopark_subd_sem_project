package main

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	db, err := pgxpool.Conn(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			User:     "admin",
			Password: "postgres",
			Port:     5432,
			Database: "forum",
		},
		MaxConnections: 50,
	})
	pgxpool
}
