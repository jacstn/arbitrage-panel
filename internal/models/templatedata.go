package models

import "github.com/jacstn/arbitrage-panel/internal/forms"

type TemplateData struct {
	Data map[string]interface{}
	Form *forms.Form
}
