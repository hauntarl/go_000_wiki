package page

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"regexp"
)

const (
	TmplView = "tmpl/view.html"
	TmplEdit = "tmpl/edit.html"
)

var (
	Templates = template.Must(template.ParseFiles(TmplView, TmplEdit))
	hyperLink = regexp.MustCompile(`\[[a-zA-Z0-9]+\]`)
)

func Load(title string) (*Page, error) {
	body, err := os.ReadFile("data/" + title + ".txt")
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

type Page struct {
	Title   string
	Body    []byte // instead of strings, as this type is expected by io libraries
	Content template.HTML
}

const rwusr = 0400 | 0200 // 0600 refer: https://linux.die.net/man/2/open

func (p *Page) Save() error {
	os.Mkdir("data", fs.ModePerm)
	return os.WriteFile("data/"+p.Title+".txt", p.Body, rwusr)
}

func (p *Page) ParseBody() {
	escape := []byte(template.HTMLEscapeString(string(p.Body)))
	p.Content = template.HTML(hyperLink.ReplaceAllFunc(escape,
		func(match []byte) []byte {
			name := match[1 : len(match)-1]
			return []byte(fmt.Sprintf(`<a href="/view/%s">%s</a>`, name, name))
		},
	))
}
