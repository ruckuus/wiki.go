package main

import (
  "fmt"
  "strings"
  "io/ioutil"
  "net/http"
)

type Page struct {
  Title string
  Body  []byte
}

var post_dir string = "_posts/"
var view_slug string = "/view/"
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
  debug("Inside indexHandler")
  files, _ := ioutil.ReadDir(post_dir)
  fmt.Fprintf(w, "<h1>Posts</h1>")
  for _, f := range files {
    post := strings.Split(f.Name(), ".")
    fmt.Fprintf(w, "<h2><a href=\"%s%s\">%s</a></h2>", view_slug, string(post[0]), string(post[0]))
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  debug("Inside viewHandler")
  title := r.URL.Path[len("/view/"):]
  p, _ := loadPage(title)

  if p != nil {
    fmt.Fprintf(w, "<h1>%s</h1><body>%s</body>", string(p.Title), string(p.Body))
  } else {
    fmt.Fprintf(w, "Internal Server Error")
  }
}

func main() {
  http.HandleFunc("/", indexHandler)
  http.HandleFunc(view_slug, viewHandler)
  http.ListenAndServe(":8080", nil)
}
