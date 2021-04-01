package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	pathRoot = "/"
	pathView = "/view/"
	pathEdit = "/edit/"
	pathSave = "/save/"
)

const (
	tmplView = "tmpl/view.html"
	tmplEdit = "tmpl/edit.html"
)

var (
	templates = template.Must(template.ParseFiles(tmplView, tmplEdit))
	validPath = regexp.MustCompile(`^/(view|edit|save)/([a-zA-Z0-9]+)$`)
	hyperLink = regexp.MustCompile(`\[[a-zA-Z0-9]+\]`)
)

func main() {
	http.HandleFunc(pathRoot, handleRoot)
	http.HandleFunc(pathView, makeHandler(handleView))
	http.HandleFunc(pathEdit, makeHandler(handleEdit))
	http.HandleFunc(pathSave, makeHandler(handleSave))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logInfo(data, message interface{}) {
	log.Printf("%-30v : %v\n", data, message)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	logInfo(r.URL.Path, "redirected to front page")
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exp := validPath.FindStringSubmatch(r.URL.Path)
		if exp == nil {
			http.NotFound(w, r)
			logInfo(r.URL.Path, "path not found")
			return
		}
		fn(w, r, exp[2])
	}
}

func handleView(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, pathEdit+title, http.StatusFound)
		logInfo(r.URL.Path, "redirected to edit page")
		return
	}

	escape := []byte(template.HTMLEscapeString(string(page.Body)))
	page.Content = template.HTML(hyperLink.ReplaceAllFunc(
		escape,
		func(match []byte) []byte {
			name := match[1 : len(match)-1]
			return []byte(fmt.Sprintf(`<a href="/view/%s">%s</a>`, name, name))
		},
	))
	render(w, "view", page)
	logInfo(page.Title, "file displayed")
}

func handleEdit(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	render(w, "edit", page)
	logInfo(page.Title, "file opened in edit mode")
}

func render(w http.ResponseWriter, tmpl string, page *Page) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func handleSave(w http.ResponseWriter, r *http.Request, title string) {
	var (
		body = r.FormValue("body")
		page = &Page{Title: title, Body: []byte(body)}
	)
	if err := page.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	http.Redirect(w, r, pathView+title, http.StatusFound)
	logInfo(page.Title, "file saved succesfully")
}

const rwusr = 0400 | 0200 // 0600 refer: https://linux.die.net/man/2/open

type Page struct {
	Title   string
	Body    []byte // instead of strings, as this type is expected by io libraries
	Content template.HTML
}

func (p *Page) save() error {
	os.Mkdir("data", fs.ModePerm)
	return os.WriteFile("data/"+p.Title+".txt", p.Body, rwusr)
}

func loadPage(title string) (*Page, error) {
	body, err := os.ReadFile("data/" + title + ".txt")
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
