package api

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/lemmego/lemmego/api/vee"
)

type TemplateData struct {
	StringMap        map[string]string
	IntMap           map[string]int
	FloatMap         map[string]float64
	BoolMap          map[string]bool
	FuncMap          template.FuncMap
	Data             map[string]any
	CSRFToken        string
	ValidationErrors vee.Errors
	Messages         []*AlertMessage
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *TemplateData) error {
	// create a template cache
	tc, err := createTemplateCache(data)
	if err != nil {
		return err
	}

	// get requested template from the cache
	t, ok := tc[tmpl]
	if !ok {
		// log the error
	}

	// render the template
	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func createTemplateCache(data *TemplateData) (map[string]*template.Template, error) {
	// create a map to act as a cache
	myCache := map[string]*template.Template{}

	// get all page files in the templates directory
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// get all partial files in the templates directory
	partials, err := filepath.Glob("./templates/*.partial.tmpl")
	if err != nil {
		return myCache, err
	}

	pages = append(pages, partials...)

	// loop through the pages one-by-one
	for _, page := range pages {
		// extract the file name (like about.page.tmpl)
		name := filepath.Base(page)

		// parse the page template file in to a template set
		ts, err := template.ParseFiles(page)
		if data.FuncMap != nil {
			ts.Funcs(data.FuncMap)
		}
		if err != nil {
			return myCache, err
		}

		// parse the layout template file in to a template set
		ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// add the template set to the cache, using the name of the page as the key
		myCache[name] = ts
	}

	// return the map
	return myCache, nil
}
