package gowiki

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// validPathExpr is a regular expression for valid paths
var validPathExpr = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// cache is a cache of template pages
var cache = template.Must(template.ParseFiles("views/edit.html", "views/view.html"))

// A Page is a page in the Wiki
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) String() string {
	return p.Title
}

func pageFilename(title string) string {
	// TODO: Build filename from title
	return fmt.Sprintf("data/%s.txt", title)
}

// Save saves the page to a file
func (p *Page) Save() error {
	filename := pageFilename(p.Title)
	return os.WriteFile(filename, p.Body, 0600)
}

// LoadPage loads a page from a file
func Load(title string) (*Page, error) {
	filename := pageFilename(title)
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPathExpr.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid page title")
	}
	return m[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	templateFile := tmpl + ".html"
	log.Printf("Render %s %s", tmpl, p.Title)
	if err := cache.ExecuteTemplate(w, templateFile, p); err != nil {
		log.Printf("Template %s failed: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := Load(title)
	if err != nil {
		log.Printf("Page not found: %s", title)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := Load(title)
	if err != nil {
		// Create new empty page
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	if err := p.Save(); err != nil {
		log.Printf("Error saving page %s: %v", title, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Saved page %s", title)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
