package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

var post_dir string = "_posts/"
var view_path string = "/view/"
var edit_path string = "/edit/"
var save_path string = "/save/"

var template_dir string = "templates/"
var templates = template.Must(template.ParseFiles(template_dir+"edit.html", template_dir+"view.html"))

var debug_enabled int = 1

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$") // 1st group: action, 2nd group: page title

func debug(msg string) {
	if debug_enabled != 0 {
		fmt.Println(msg)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	debug("Inside getTitle with parameter:")
	m := validPath.FindStringSubmatch(r.URL.Path)
	debug(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}

	return m[2], nil // Title is in the 2nd group.
}

func (p *Page) save() error {
	filename := post_dir + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := post_dir + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, p *Page, template_name string) {

	err := templates.ExecuteTemplate(w, template_name+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	debug("Inside indexHandler")
	files, _ := ioutil.ReadDir(post_dir)
	fmt.Fprintf(w, "<h1>Posts</h1>")
	for _, f := range files {
		post := strings.Split(f.Name(), ".")
		fmt.Fprintf(w, "<h2><a href=\"%s%s\">%s</a></h2>", view_path, string(post[0]), string(post[0]))
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	debug("Inside viewHandler")
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, p, "view")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	debug("Inside editHandler")
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	p, err := loadPage(title)

	if err != nil {
		body := []byte("New Article")
		if p != nil {
			body = p.Body
		}
		p = &Page{Title: title, Body: body}
	}
	renderTemplate(w, p, "edit")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	debug("Inside saveHandler")
	title, err := getTitle(w, r)
	if err != nil {
		return
	}

	body := r.FormValue("body")

	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, view_path+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc(view_path, viewHandler)
	http.HandleFunc(edit_path, editHandler)
	http.HandleFunc(save_path, saveHandler)
	http.ListenAndServe(":8080", nil)
}
