package handlers

import (
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

var app *config.AppConfig

func NewHandlers(c *config.AppConfig) {
	app = c
}

func About(w http.ResponseWriter, r *http.Request) {
	//ip := app.Session.GetString(r.Context(), "remote_ip")
	data := models.TemplateData{}

	renderTemplate(w, "about", &data)
}

func RunningTrades(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["csrf_token"] = nosurf.Token(r)

	data["trade_list"] = models.ListRunningTrades(app.DB)

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
		fmt.Println()
		name := f.Name()
		if strings.Contains(name, ".csv") {
			fl = append(fl, f.Name())
		}
	}

	render.JSON(w, r, fl)
}

func renderTemplate(w http.ResponseWriter, templateName string, data *models.TemplateData) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+templateName+".go.tmpl", "./templates/base.layout.go.tmpl")

	err := parsedTemplate.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "Error handling template page!!", err)
	}
}
