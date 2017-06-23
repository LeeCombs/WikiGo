package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"errors"
)

type Page struct {
	Title string
	Body []byte
}

var (
	cwd, _ = os.Getwd()
	templatePath = filepath.Join(cwd, "/src/templates/")
	viewPath = filepath.Join(templatePath, "/view/")
	editPath = filepath.Join(templatePath, "/edit/")
	savePath = filepath.Join(templatePath, "/save/")
	templates = template.Must(template.ParseFiles(editPath + ".html", viewPath + ".html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

/* Validate path and extract the Page title */
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
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
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/*  */
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

/* Handle generic requests */
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/* Handle URLs prefixed with /view/. Allows user to view a Page */
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
    	// If the page wasn't found, redirect to the edit view
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

/* Handle URLs prefixed with /edit/. Allows user to edit a Page */
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

/* Handle URLs prefixed with /save/. Allows user to save a Page */
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
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
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8282", nil)
}