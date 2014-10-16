package main

import (
  "fmt"
  "strings"
  "io/ioutil"
  "net/http"
  "html/template"
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

var debug_enabled int = 1

func debug(msg string) {
  if debug_enabled != 0 {
    fmt.Println(msg)
  }
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
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
    template_path := template_dir + template_name + ".html"

    t, _ := template.ParseFiles(template_path)
    t.Execute(w, p)
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
  title := r.URL.Path[len(view_path):]
  p, _ := loadPage(title)

  if p != nil {
    renderTemplate(w, p, "view")
  } else {
    fmt.Fprintf(w, "Internal Server Error")
  }
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  debug("Inside editHandler")
  title := r.URL.Path[len(edit_path):]
  p, _ := loadPage(title)

  if p != nil {
    renderTemplate(w, p, "edit")
  } else {
    fmt.Fprintf(w, "Internal Server Error")
  }
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc(view_path, viewHandler)
  http.HandleFunc(edit_path, editHandler)
  http.HandleFunc(save_path, saveHandler)
  http.ListenAndServe(":8080", nil)
}
