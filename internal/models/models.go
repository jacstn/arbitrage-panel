package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Url struct {
	Id        uint16
	Name      string
	CreatedAt string
}

func SaveUrl(db *sql.DB, url Url) int64 {
	var id int64
	err := db.QueryRow("SELECT id FROM domain where name=?", url.Name).Scan(&id)

	if err == nil && id > 0 {
		return id
	}

	res, err := db.Exec("Insert into domain (name, createdAt) values (?, ?)", url.Name, time.Now())

	if err != nil {
		fmt.Println("error while inserting into database", err)
		return 0
	}

	id, err = res.LastInsertId()
	if err != nil {
		println("Error:", err.Error())
		return 0
	}

	fmt.Println("record saved, id", id)
	return id
}

func ListUrls(db *sql.DB) []Url {
	urls := []Url{}

	res, err := db.Query("SELECT * FROM domain")

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Url{}
	}

	for res.Next() {

		var url Url
		err := res.Scan(&url.Id, &url.Name, &url.CreatedAt)

		if err != nil {
			fmt.Println(err)
		}
		urls = append(urls, url)
	}

	return urls
}

func GetUrlById(db *sql.DB, id string) Url {
	url := Url{}
	err := db.QueryRow("SELECT * FROM domain where id=?", id).Scan(&url.Id, &url.Name, &url.CreatedAt)

	if err != nil {

		return url
	}

	return url
}
