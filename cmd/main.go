package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
)

func main() {
	dbConn := pgx.Conn{
		User:                 "farcoad",
		Database:             "forum",
		Password:             "postgres",
		PreferSimpleProtocol: false,
	}
	db, err := pgxpool.ConnectConfig(context.Background(),
	{  User: "admin",
			Password: "postgres",
		Port:     5432,
		Database: "forum",

	})
	if err != nil {
		fmt.Print(err)
	}


	dbConf := pgx.ConnConfig{
		User:                 "farcoad",
		Database:             "forum",
		Password:             "postgres",
		PreferSimpleProtocol: false,
	}

	pgx.
}
