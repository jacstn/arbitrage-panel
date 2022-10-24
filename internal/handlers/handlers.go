package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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

func Home(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["csrf_token"] = nosurf.Token(r)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	per_page := 30

	lt, nt := models.ListTrades(app.DB, page, int(per_page))
	data["trade_list"] = lt
	data["num_of_trades"] = nt
	data["page"] = page
	fmt.Println(data["page"])
	renderTemplate(w, "home", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func ViewUrl(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	id := app.Session.Pop(r.Context(), "saved_id")

	fmt.Print(id)

	renderTemplate(w, "view-url", &models.TemplateData{
		Data: data,
	})
}

func renderTemplate(w http.ResponseWriter, templateName string, data *models.TemplateData) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+templateName+".go.tmpl", "./templates/base.layout.go.tmpl")

	err := parsedTemplate.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Error handling template page!!", err)
	}
}
