package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jacstn/arbitrage-panel/config"
	"github.com/jacstn/arbitrage-panel/internal/forms"
	"github.com/jacstn/arbitrage-panel/internal/models"
	"github.com/justinas/nosurf"
)

const (
	OpenLong   = "long from binance"
	CloseLong  = "SPOT closed long info"
	OpenShort  = "CROSS short from binance"
	CloseShort = "CROSS closed short info"
	IncrShort  = "increasing CROSS short from binance"
	IncrLong   = "increasing long from binance"
)

var app *config.AppConfig

func NewHandlers(c *config.AppConfig) {
	app = c
}

func About(w http.ResponseWriter, r *http.Request) {
	//ip := app.Session.GetString(r.Context(), "remote_ip")
	data := models.TemplateData{}

	renderTemplate(w, "about", &data)
}

type BinanceTransaction struct {
	CummulativeQuoteQty string                 `json:"cummulativeQuoteQty"`
	X                   map[string]interface{} `json:"-"`
}

func getTradeResult(tid uint64) (float32, int) {
	raws := models.GetLogs(app.DB, tid)
	var res float32 = 0.0
	numOfInc := 0
	for _, v := range raws {
		if v.Message == OpenShort || v.Message == CloseLong || v.Message == IncrShort {
			trans := BinanceTransaction{}
			if err := json.Unmarshal([]byte(v.Raw), &trans); err != nil {
				fmt.Println("error while unmarshaling json", err)
			} else if s, err := strconv.ParseFloat(trans.CummulativeQuoteQty, 64); err == nil {
				res -= float32(s)
			} else {
				fmt.Println("error while float conv", err)
			}
		}

		if v.Message == OpenLong || v.Message == CloseShort || v.Message == IncrLong {
			trans := BinanceTransaction{}
			if err := json.Unmarshal([]byte(v.Raw), &trans); err != nil {
				fmt.Println("error while un-marshaling json", err)
			} else if s, err := strconv.ParseFloat(trans.CummulativeQuoteQty, 64); err == nil {
				res += float32(s)
			} else {
				fmt.Println("error while float conv", err)
			}
		}

		if v.Message == IncrLong {
			numOfInc++
		}
	}

	return res, numOfInc
}

func RunningTrades(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["csrf_token"] = nosurf.Token(r)
	btcPrice := models.GetPrice(app.DB, "BTCUSDT")
	rt := models.ListRunningTrades(app.DB)

	var usdt_res float32 = 0.0
	var btc_res float32 = 0.0

	for i, v := range rt {
		res, noi := getTradeResult(v.Id)
		res = res - v.ValShort + v.ValLong
		if rt[i].SymbolShort[len(rt[i].SymbolShort)-3:] == "BTC" {
			btc_res += res
			rt[i].CurrRes = res
			rt[i].CurrResDisp = fmt.Sprintf("%.8f (%.2f)", btc_res, btc_res*btcPrice)
		} else {
			usdt_res += res
			rt[i].CurrRes = res
			rt[i].CurrResDisp = fmt.Sprintf("%.2f", res)
		}

		rt[i].IncNo = noi
	}

	data["trade_list"] = rt
	data["btc_res"] = btc_res
	data["usdt_res"] = usdt_res
	data["btc_res_usdt"] = btcPrice * btc_res
	renderTemplate(w, "running-trades", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func AllTrades(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["csrf_token"] = nosurf.Token(r)

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	searchText := r.URL.Query().Get("search")
	statuses := r.URL.Query().Get("status")

	fmt.Println("search text from url", statuses)

	per_page := 200

	lt, nt := models.ListTrades(app.DB, searchText, statuses, page, int(per_page))

	data["trade_list"] = lt
	data["num_of_trades"] = nt
	data["page"] = page
	data["searchText"] = searchText

	renderTemplate(w, "all-trades", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func GetLogs(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	data := models.GetLogs(app.DB, uint64(id))

	render.JSON(w, r, data)
}

func GetTradeById(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	data["trade"] = models.GetTradeById(app.DB, uint64(id))

	renderTemplate(w, "trade_detail", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func Market(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	renderTemplate(w, "market", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	pyPath := os.Getenv("ARBITRAGE_PY_DIR")
	if pyPath == "" {
		http.Error(w, pyPath, 400)
		return
	}

	files, err := ioutil.ReadDir(pyPath + "/data")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "cannot read python files path: "+pyPath, 400)
		return
	}
	var fl []string
	for _, f := range files {
		name := f.Name()
		if strings.Contains(name, ".csv") {
			fl = append(fl, f.Name())
		}
	}

	render.JSON(w, r, fl)
}

func renderTemplate(w http.ResponseWriter, templateName string, data *models.TemplateData) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+templateName+".go.html", "./templates/base.layout.go.html")

	err := parsedTemplate.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "Error handling template page!!", err)
	}
}
