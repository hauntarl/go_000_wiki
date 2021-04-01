package page

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"regexp"
)

var (
	// Parse templates files only once, re-use them as and when required
	Templates = template.Must(template.ParseFiles("tmpl/view.html", "tmpl/edit.html"))
	// to capture internal-linking
	hyperLink = regexp.MustCompile(`\[[a-zA-Z0-9]+\]`)
)

// Loads the given title from data directory and return Page structure
func Load(title string) (*Page, error) {
	body, err := os.ReadFile("data/" + title + ".txt")
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

type Page struct {
	Title   string        // name of file, without extension
	Body    []byte        // instead of strings, as this type is expected by io libraries
	Content template.HTML // ExecuteTemplate will not escape variables of the type template.HTML
}

const rwusr = 0400 | 0200 // 0600 refer: https://linux.die.net/man/2/open

// Persist current Page structure into database
func (p *Page) Save() error {
	os.Mkdir("data", fs.ModePerm)
	return os.WriteFile("data/"+p.Title+".txt", p.Body, rwusr)
}

// Populates Content attribute of Page structure
//
// 1. searches for words marked to be inter-linked
//
// 2. replaces the matches with appropriate href syntax for html
func (p *Page) ParseBody() {
	escape := []byte(template.HTMLEscapeString(string(p.Body)))
	p.Content = template.HTML(hyperLink.ReplaceAllFunc(escape,
		func(match []byte) []byte {
			name := match[1 : len(match)-1]
			return []byte(fmt.Sprintf(`<a href="/view/%s">%s</a>`, name, name))
		},
	))
}
