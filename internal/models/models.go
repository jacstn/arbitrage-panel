package models

import (
	"database/sql"
	"fmt"
)

type Trade struct {
	Id        uint16
	Name      string
	CreatedAt string
}

func ListTrades(db *sql.DB) []Trade {
	trades := []Trade{}

	res, err := db.Query("SELECT * FROM domain")

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}
	}

	for res.Next() {

		var trade Trade
		err := res.Scan(&trade.Id, &trade.Name, &trade.CreatedAt)

		if err != nil {
			fmt.Println(err)
		}
		trades = append(trades, trade)
	}

	return trades
}

func GetLogs(db *sql.DB, id string) Trade {
	logs := Trade{}
	err := db.QueryRow("SELECT * FROM domain where id=?", id).Scan(&logs.Id, &logs.Name, &logs.CreatedAt)

	if err != nil {

		return logs
	}

	return logs
}
