package models

import (
	"database/sql"
	"fmt"
)

type Trade struct {
	Id        uint16
	LongQty   float32
	ShortQty  float32
	OpenedAt  string
	UpdatedAt string
}

type TradeLogs struct {
	Id       uint16
	Category string
	Message  string
}

func ListTrades(db *sql.DB) []Trade {
	trades := []Trade{}

	res, err := db.Query("SELECT * FROM trades")

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}
	}

	for res.Next() {

		var trade Trade
		err := res.Scan(&trade.Id, &trade.LongQty, &trade.UpdatedAt)

		if err != nil {
			fmt.Println(err)
		}
		trades = append(trades, trade)
	}

	return trades
}

func GetLogs(db *sql.DB, tradeId uint32) TradeLogs {
	logs := TradeLogs{}
	err := db.QueryRow("SELECT id, category, message FROM trade_logs where trade_id=?", tradeId).Scan(&logs.Id, &logs.Category, &logs.Message)

	if err != nil {

		return logs
	}

	return logs
}
