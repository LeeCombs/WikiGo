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

/* Handle generic requests */
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/* Handle URLs prefixed with /view/. Allows user to view a Page */
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

/* Handle URLs prefixed with /edit/. Allows user to edit a Page */
func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }

	cwd, _ := os.Getwd()
    t, err := template.ParseFiles(filepath.Join(cwd, "/src/templates/edit.html"))
    if err != nil {
        fmt.Println(err)
        return
    }
    t.Execute(w, p)
}

/* Handle URLs prefixed with /save/. Allows user to save a Page */
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8282", nil)
}