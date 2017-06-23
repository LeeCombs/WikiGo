package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"os"
	"path/filepath"
)

type Page struct {
	Title string
	Body []byte
}

/* Save a Page locally */
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/* Load a locally saved Page */
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

/* Render HTML templates */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	cwd, err := os.Getwd()
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    t, err := template.ParseFiles(filepath.Join(cwd, "/src/templates/" + tmpl + ".html"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

/* Handle generic requests */
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/* Handle URLs prefixed with /view/. Allows user to view a Page */
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)

    // If the page wasn't found, redirect to the edit view
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }

    renderTemplate(w, "view", p)
}

/* Handle URLs prefixed with /edit/. Allows user to edit a Page */
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }

    renderTemplate(w, "edit", p)
}

/* Handle URLs prefixed with /save/. Allows user to save a Page */
func saveHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}

    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // After saving, send the user to the view page
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8282", nil)
}