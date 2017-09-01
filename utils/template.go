package utils

import (
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/entities"
	"html/template"
	"net/http"
	"path/filepath"
)

// templatePrefix is the constant containing the filepath prefix for templates.
const templatePrefix = "views"

// baseTemplates is the array of 'base' templates used in each template render.
var baseTemplates = []string{
	"partials/header.tmpl",
	"partials/footer.tmpl",
	"partials/base.tmpl",
}

type Template struct {
	config *entities.Config
}

func NewTemplate(config *entities.Config) *Template {
	return &Template{
		config,
	}
}

// Render renders a template on the ResponseWriter w.
//
// The interface{} data gives the data to be sent to the template.
//
// The string mainTmpl gives the name, relative to views, of the main
// template to render.  The variadic argument addTmpls names any additional
// templates mainTmpl depends on.
//
// Render returns any error that occurred when rendering the template.
func (t *Template) Render(w http.ResponseWriter, data interface{}, mainTemplate string, additionalTemplates ...string) error {

	var err error

	templates := append(baseTemplates, additionalTemplates...)

	var extraTemplateAssets []string

	for _, tmpl := range templates {
		ta := MustAsset(filepath.Join(templatePrefix, tmpl))
		extraTemplateAssets = append(extraTemplateAssets, string(ta))
	}

	mainTemplateAsset := string(MustAsset(filepath.Join(templatePrefix, mainTemplate)))

	tempFuncs := template.FuncMap{
		"html":  renderHTML,
		"url":   func() string { return t.config.RootUrl },
		"proto": func() string { return t.config.Protocol },
	}

	temp, err := parseTemplates(mainTemplateAsset, tempFuncs, extraTemplateAssets)

	if err != nil {
		return err
	}

	return temp.Execute(w, data)

}

func parseTemplates(mainTemplate string, functions template.FuncMap, additionalTemplates []string) (*template.Template, error) {
	t := template.New("_all").Funcs(functions)
	var err error
	for i, tmpl := range additionalTemplates {
		t, err = t.New(fmt.Sprint("_", i)).Parse(tmpl)
		if err != nil {
			return nil, err
		}
	}
	t, err = t.Parse(mainTemplate)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// renderHTML takes some html as a string and returns a template.HTML
//
// Handles plain text gracefully.
func renderHTML(value interface{}) template.HTML {
	return template.HTML(fmt.Sprint(value))
}
