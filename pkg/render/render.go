package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/msaufi2325/06_bookings/pkg/config"
	"github.com/msaufi2325/06_bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates set the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate render templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		fmt.Println("error in tc[tmpl]")
		fmt.Println(tmpl)
		fmt.Println(tc)
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		fmt.Println("error in t.Execute")
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error in buf.WriteTo")
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.gohtml from the ./templates directory
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		fmt.Println("error in filepath.Glob")
		return myCache, err
	}

	// range through all files ending with *.page.gohtml
	for _, page := range pages {
		// get the file name
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			fmt.Println("error in template.New")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*layout.gohtml")
		if err != nil {
			fmt.Println("error in filepath.Glob")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				fmt.Println("error in ts.ParseGlob")
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil

}
