package models

import (
	"database/sql"
	"fmt"
)

type Trade struct {
	Id          uint64
	Status      string
	SymbolLong  string
	SymbolShort string
	TimeOrigin  string
	OpenDiff    float32
	QtyLong     float32
	QtyShort    float32
	OpenedAt    string
	UpdatedAt   string
}

type TradeLogs struct {
	Id       uint16
	Category string
	Message  string
}

func ListTrades(db *sql.DB) ([]Trade, uint64) {
	trades := []Trade{}

	res, err := db.Query("SELECT count(id) FROM trades order by openedAt DESC  LIMIT 5")

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}, 0
	}
	var totalNum uint64
	res.Scan(&totalNum)

	res, err = db.Query("SELECT id, status, symbol_long, symbol_short, time_origin, open_diff, qty_long, qty_short, openedAt, updatedAt FROM trades order by openedAt DESC  LIMIT 5")

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}, 0
	}

	for res.Next() {

		var trade Trade
		err := res.Scan(&trade.Id, &trade.Status, &trade.SymbolLong, &trade.SymbolShort, &trade.TimeOrigin, &trade.OpenDiff, &trade.QtyLong, &trade.QtyShort, &trade.OpenedAt, &trade.UpdatedAt)

		if err != nil {
			fmt.Println(err)
		}
		trades = append(trades, trade)
	}

	return trades, totalNum
}

func GetLogs(db *sql.DB, tradeId uint32) TradeLogs {
	logs := TradeLogs{}
	err := db.QueryRow("SELECT id, category, message FROM trade_logs where trade_id=?", tradeId).Scan(&logs.Id, &logs.Category, &logs.Message)

	if err != nil {

		return logs
	}

	return logs
}
