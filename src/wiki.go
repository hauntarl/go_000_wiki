package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/go_000_wiki/page"
)

const (
	pathRoot = "/"
	pathHome = "/view/FrontPage"
	pathView = "/view/"
	pathEdit = "/edit/"
	pathSave = "/save/"
)

func main() {
	http.HandleFunc(pathRoot, handleRoot)
	http.HandleFunc(pathView, makeHandler(handleView))
	http.HandleFunc(pathEdit, makeHandler(handleEdit))
	http.HandleFunc(pathSave, makeHandler(handleSave))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, pathHome, http.StatusFound)
	logInfo(r.URL.Path, "redirected to front page")
}

// acceptable combinations of path by the application
var validPath = regexp.MustCompile(`^/(view|edit|save)/([a-zA-Z0-9]+)$`)

// wrapper function, which will handle extraction of title, required by all
func makeHandler(handler func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exp := validPath.FindStringSubmatch(r.URL.Path)
		if exp == nil {
			http.NotFound(w, r)
			logInfo(r.URL.Path, "path not found")
			return
		}
		handler(w, r, exp[2])
	}
}

// displays content for given title, if absent, opens the file in edit mode
func handleView(w http.ResponseWriter, r *http.Request, title string) {
	p, err := page.Load(title)
	if err != nil {
		http.Redirect(w, r, pathEdit+title, http.StatusFound)
		logInfo(r.URL.Path, "redirected to edit page")
		return
	}

	p.ParseBody()
	render(w, "view", p)
	logInfo(p.Title, "file displayed")
}

// loads the title to be edited from database
func handleEdit(w http.ResponseWriter, r *http.Request, title string) {
	p, err := page.Load(title)
	if err != nil {
		p = &page.Page{Title: title}
	}
	render(w, "edit", p)
	logInfo(p.Title, "file opened in edit mode")
}

// helper function which handles the rendering template with given Page data
func render(w http.ResponseWriter, tmpl string, p *page.Page) {
	if err := page.Templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

// Saves created/edited Page into database
func handleSave(w http.ResponseWriter, r *http.Request, title string) {
	var (
		body = r.FormValue("body")
		p    = &page.Page{Title: title, Body: []byte(body)}
	)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	http.Redirect(w, r, pathView+title, http.StatusFound)
	logInfo(p.Title, "file saved succesfully")
}

func logInfo(data, message interface{}) {
	log.Printf("%-30v : %v\n", data, message)
}
