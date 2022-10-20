package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	var db_user string
	if os.Getenv("DB_USER") == "" {
		db_user = "root"
	} else {
		db_user = os.Getenv("DB_USER")
	}

	var db_pass string
	if os.Getenv("DB_PASS") == "" {
		db_pass = ""
	} else {
		db_pass = os.Getenv("DB_PASS")
	}

	var db_name string
	if os.Getenv("DB_NAME") == "" {
		db_name = "arbitrage"
	} else {
		db_name = os.Getenv("DB_NAME")
	}

	db, err := sql.Open("mysql", db_user+":"+db_pass+"@tcp(127.0.0.1:3306)/"+db_name)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Succesfully connected to Database")

	return db
}
