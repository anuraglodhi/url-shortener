package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/AnuragLodhi/urlshortener/database"
	"github.com/lucsky/cuid"
)

var db database.Database

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := strings.Split(r.URL.Path, "/")[0]

	longUrl, err := db.GetLongUrl(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	originalurl, err := url.Parse(longUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if originalurl.Scheme == "" {
		originalurl.Scheme = "http"
	}
	http.Redirect(w, r, originalurl.String(), http.StatusSeeOther)
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
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

	short, err := db.GetShortUrl(url)
	if err != nil {
		for {
			short = cuid.Slug()
			exists, err := db.ShortUrlExists(short)
			if err != nil {
				http.Error(w, "Link could not be generated. Try again", http.StatusInternalServerError)
			}
			if !exists {
				break
			}
		}
		err = db.InsertUrl(short, url)
		if err != nil {
			http.Error(w, "Link could not be generated. Try again", http.StatusInternalServerError)
		}
	}

	link := "http://localhost:8080/" + short
	tmpl, err := template.New("shortlink").Parse("Short link: <a href='" + link + "'>" + link + "</>")
	if err != nil {
		http.Error(w, "Link could not be generated. Try again", http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)
}

func main() {
	db = database.New()

	r := http.NewServeMux()

	r.HandleFunc("/", handleIndexPage)
	r.HandleFunc("/shorten", handleShorten)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", r))
}
