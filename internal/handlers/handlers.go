package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jacstn/arbitrage-panel/config"
	"github.com/jacstn/arbitrage-panel/internal/forms"
	"github.com/jacstn/arbitrage-panel/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewHandlers(c *config.AppConfig) {
	app = c
}

func List(w http.ResponseWriter, r *http.Request) {

	displayData := make(map[string]interface{})

	app.Session.Put(r.Context(), "remote_ip", r.Host)
	displayData["list_of_trades"] = models.ListTrades(app.DB)

	renderTemplate(w, "home", &models.TemplateData{
		Data: displayData,
	})
}

func About(w http.ResponseWriter, r *http.Request) {
	//ip := app.Session.GetString(r.Context(), "remote_ip")
	data := models.TemplateData{}

	renderTemplate(w, "about", &data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["url_model"] = models.Trade{}
	data["csrf_token"] = nosurf.Token(r)

	renderTemplate(w, "new-url", &models.TemplateData{
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
