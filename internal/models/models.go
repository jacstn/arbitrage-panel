package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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

type TradeLog struct {
	Id       uint64
	Category string
	Message  string
	Raw      string
}

func ListTrades(db *sql.DB, searchText string, page int, perPage int) ([]Trade, uint64) {
	trades := []Trade{}
	log.Println(searchText)
	where := "WHERE `status` not in ('NO_BALANCE', 'CANNOT_BORROW') "

	searchText = strings.ToLower(searchText)

	if searchText != "" {
		where += fmt.Sprintf("and LOWER(symbol_long) LIKE '%%%s%%' OR LOWER(symbol_short) LIKE '%%%s%%' OR LOWER(status) LIKE '%%%s%%'", searchText, searchText, searchText)
	}

	res, err := db.Query(fmt.Sprintf("SELECT count(id) FROM trades %s", where))

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}, 0
	}
	var totalNum uint64
	res.Scan(&totalNum)

	res, err = db.Query(fmt.Sprintf("SELECT id, status, symbol_long, symbol_short, time_origin, open_diff, qty_long, qty_short, openedAt, updatedAt FROM trades %s ORDER BY openedAt DESC  LIMIT %d OFFSET %d", where, perPage, (page-1)*perPage))

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

func GetLogs(db *sql.DB, tradeId uint64) []TradeLog {
	logs := []TradeLog{}
	res, err := db.Query("SELECT id, category, message, raw FROM trade_logs where trade_id=?", tradeId)

	if err != nil {
		return logs
	}

	for res.Next() {
		var log TradeLog
		var raw sql.NullString
		err := res.Scan(&log.Id, &log.Category, &log.Message, &raw)
		log.Raw = raw.String
		if err != nil {
			fmt.Println(err)
		}
		logs = append(logs, log)
	}

	return logs
}
