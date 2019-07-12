package api

import (
	"os"
	"log"
	"path/filepath"
	"net/http"
	"text/template"
)

var templates *template.Template

func init() {
	cwd, _ := os.Getwd()
	tmpl := filepath.Join(cwd, "./templates/index.html")	
	templates = template.Must(template.ParseFiles(tmpl))
}

// Index handle index.html page
func index(w http.ResponseWriter, r *http.Request) {	
 	if r.Method == "GET" {
 		templates.Execute(w, nil)
 	}
 }
