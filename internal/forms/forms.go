package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) ValidUrl(field string, r *http.Request) bool {
	_, err := url.ParseRequestURI(r.Form.Get(field))
	if err != nil {
		f.Errors.Add(field, "This URL is invalid")
		return false
	}

	return true
}

func (f *Form) Has(field string, r *http.Request) bool {
	if r.Form.Get(field) == "" {
		f.Errors.Add(field, "This field cannot be empty")
		return false
	}

	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
