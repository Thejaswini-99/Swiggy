package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {

	var err error

	db, err = sql.Open("mysql", "root:12345@tcp(localhost:3306)/swiggy")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB Connected Successfully")
}
