package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Trade struct {
	Id           uint64
	Status       string
	SymbolLong   string
	SymbolShort  string
	TimeOrigin   string
	OpenDiff     float32
	QtyLong      float32
	QtyShort     float32
	OpenedAt     string
	OpenedAgo    string
	HoursToClose string
	UpdatedAt    string
	ValLong      float32
	ValShort     float32
	CurrRes      string
}

type TradeLog struct {
	Id       uint64
	Ago      string
	Category string
	Message  string
	Raw      string
}

type BnbTrans struct {
	Id        uint64
	Raw       string
	Price     float64
	Qty       float64
	CreatedAt string
	Mode      string
}

func ListRunningTrades(db *sql.DB) []Trade {
	trades := []Trade{}

	res, err := db.Query(`SELECT (select price from prices where symbol=symbol_long ORDER BY time DESC LIMIT 1) * qty_long as val_long, 
	(select price from prices where symbol=symbol_short ORDER BY time DESC LIMIT 1) * qty_short as val_short, 
	id, status, symbol_long, symbol_short, time_origin, open_diff, qty_long, qty_short, openedAt, 
	TIMEDIFF(NOW(), openedAt) AS opened_ago, hours_to_close, updatedAt FROM trades where status in ('RUNNING', 'MANUAL') ORDER BY openedAt DESC`)

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}
	}

	for res.Next() {
		var trade Trade
		err := res.Scan(&trade.ValLong, &trade.ValShort,
			&trade.Id, &trade.Status, &trade.SymbolLong,
			&trade.SymbolShort, &trade.TimeOrigin, &trade.OpenDiff,
			&trade.QtyLong, &trade.QtyShort, &trade.OpenedAt,
			&trade.OpenedAgo, &trade.HoursToClose,
			&trade.UpdatedAt)

		if err != nil {
			fmt.Println(err)
		}
		trades = append(trades, trade)
	}

	return trades
}

func GetPrice(db *sql.DB, symbol string) float32 {
	var price float32

	err := db.QueryRow(fmt.Sprintf("SELECT price from prices where symbol='%s' order by time desc limit 1", symbol)).Scan(&price)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return price
}

func ListTrades(db *sql.DB, searchText string, status string, page int, perPage int) ([]Trade, uint64) {
	trades := []Trade{}
	log.Println(searchText)
	where := ""
	if status != "" {
		where += fmt.Sprintf("WHERE `status` in (%s)", status)
	}

	searchText = strings.ToLower(searchText)

	if searchText != "" {
		if where == "" {
			where += "WHERE"
		} else {
			where += " AND "
		}
		where += fmt.Sprintf(" LOWER(symbol_long) LIKE '%%%s%%' OR LOWER(symbol_short) LIKE '%%%s%%' OR LOWER(status) LIKE '%%%s%%' OR CAST(id as CHAR) LIKE '%%%s%%'", searchText, searchText, searchText, searchText)
	}
	query := fmt.Sprintf("SELECT count(id) FROM trades %s", where)
	fmt.Println("query", query)
	res, err := db.Query(query)

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}, 0
	}
	var totalNum uint64
	res.Scan(&totalNum)

	res, err = db.Query(fmt.Sprintf("SELECT id, status, symbol_long, symbol_short, time_origin, open_diff, qty_long, qty_short, openedAt, TIMEDIFF(NOW(), openedAt) AS opened_ago, updatedAt FROM trades %s ORDER BY openedAt DESC  LIMIT %d OFFSET %d", where, perPage, (page-1)*perPage))

	if err != nil {
		fmt.Println("cannot query from database", err)
		return []Trade{}, 0
	}

	for res.Next() {

		var trade Trade
		err := res.Scan(&trade.Id, &trade.Status, &trade.SymbolLong, &trade.SymbolShort, &trade.TimeOrigin, &trade.OpenDiff, &trade.QtyLong, &trade.QtyShort, &trade.OpenedAt, &trade.OpenedAgo, &trade.UpdatedAt)

		if err != nil {
			fmt.Println(err)
		}
		trades = append(trades, trade)
	}

	return trades, totalNum
}

func GetLogs(db *sql.DB, tradeId uint64) []TradeLog {
	logs := []TradeLog{}
	res, err := db.Query("SELECT id, category, message, raw, TIMEDIFF(NOW(), createdAt) AS createdAgo FROM trade_logs where trade_id=?", tradeId)

	if err != nil {
		return logs
	}

	for res.Next() {
		var log TradeLog
		var raw sql.NullString
		err := res.Scan(&log.Id, &log.Category, &log.Message, &raw, &log.Ago)
		log.Raw = raw.String
		if err != nil {
			fmt.Println(err)
		}

		logs = append(logs, log)
	}

	return logs
}

func GetBinanceTransactionsFromLogs(db *sql.DB, tradeId uint64) []TradeLog {
	logs := []TradeLog{}
	res, err := db.Query("SELECT raw FROM trade_logs where trade_id=? AND message in ('CROSS short from binance', 'long from binance', 'increasing CROSS short from binance','increasing long from binance')", tradeId)

	if err != nil {
		return logs
	}

	for res.Next() {
		var log TradeLog
		var raw sql.NullString
		err := res.Scan(&log.Id, &log.Category, &log.Message, &raw, &log.Ago)
		log.Raw = raw.String
		if err != nil {
			fmt.Println(err)
		}

		logs = append(logs, log)
	}

	return logs
}

func GetTradeById(db *sql.DB, id uint64) Trade {
	var trade Trade

	err := db.QueryRow(fmt.Sprintf("SELECT id, status, symbol_long, symbol_short, time_origin, open_diff, qty_long, qty_short, openedAt, TIMEDIFF(NOW(), openedAt) AS opened_ago, updatedAt FROM trades where id=%d", id)).Scan(&trade.Id, &trade.Status, &trade.SymbolLong, &trade.SymbolShort, &trade.TimeOrigin, &trade.OpenDiff, &trade.QtyLong, &trade.QtyShort, &trade.OpenedAt, &trade.OpenedAgo, &trade.UpdatedAt)

	if err != nil {
		fmt.Println(err)
		return Trade{}
	}

	return trade
}

func GetListOfStatuses(db *sql.DB) []string {
	statuses := []string{}
	res, err := db.Query("SELECT distinct(status) from trades")

	if err != nil {
		return statuses
	}

	for res.Next() {
		var st string
		err := res.Scan(&st)

		if err != nil {
			fmt.Println(err)
		}

		statuses = append(statuses, st)
	}

	return statuses
}

func GetListOfBnbTransactions(db *sql.DB) []BnbTrans {
	res, err := db.Query("SELECT id, raw, createdAt, price, qty, mode from bnb_trans")
	var bnbTrans []BnbTrans
	if err != nil {
		return []BnbTrans{}
	}

	for res.Next() {
		var bt BnbTrans
		err := res.Scan(&bt.Id, &bt.Raw, &bt.CreatedAt, &bt.Price, &bt.Qty, &bt.Mode)

		if err != nil {
			fmt.Println(err)
		}

		bnbTrans = append(bnbTrans, bt)
	}

	return bnbTrans
}
