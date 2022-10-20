package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db_user := os.Getenv("DB_USER")
	if db_user == "" {
		db_user = "root"
	}

	db_pass := os.Getenv("DB_PASS")

	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		db_name = "arbitrage"
	}

	db_url := db_user + ":" + db_pass + "@tcp(127.0.0.1:3306)/" + db_name
	fmt.Println(db_url)
	db, err := sql.Open("mysql", db_url)
	if err != nil {
		fmt.Println("Unable to connect to mysql server using URL", db_url)
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Succesfully connected to Database")

	return db
}
