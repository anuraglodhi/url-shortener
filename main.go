package main

import (
	// "fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

var urls map[string]string = make(map[string]string)

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := strings.Split(r.URL.Path, "/")[0]
	
	originalString, ok := urls[short]
	originalurl, err := url.Parse(originalString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	if ok {
		if originalurl.Scheme == "" {
			originalurl.Scheme = "http";
		}
		http.Redirect(w, r, originalurl.String(), http.StatusSeeOther)
	} else { 
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/")
	if url != "" {
		r.URL.Path = url
		handleRedirect(w, r)
		return
	}
	tmpl, err := template.ParseFiles("pages/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	short := url[:3]
	urls[short] = url
	io.WriteString(w, url + " shortened to " + short)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", handleShorten)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
