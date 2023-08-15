package gowiki

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("views/" + tmpl + ".html")
	if err != nil {
		log.Printf("Template %s failed: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Render %s %s", tmpl, p.Title)
	if err := t.Execute(w, p); err != nil {
		log.Printf("Template %s failed: %v", tmpl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := Load(title)
	if err != nil {
		log.Printf("Page not found: %s", title)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := Load(title)
	if err != nil {
		// Create new empty page
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
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
