package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/lucsky/cuid"
)

var db *sql.DB

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := strings.Split(r.URL.Path, "/")[0]

	var originalString string
	row := db.QueryRow(fmt.Sprintf("select longurl from urls where shorturl=\"%s\"", short))
	err := row.Scan(&originalString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	originalurl, err := url.Parse(originalString)
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
	var short string

	row := db.QueryRow(fmt.Sprintf("select shorturl from urls where longurl=\"%s\"", url))
	err := row.Scan(&short)
	if err != nil {
		for {
			short = cuid.Slug()
			row := db.QueryRow(fmt.Sprintf("select count(*) from urls where shorturl=\"%s\"", short))
			var count int
			row.Scan(&count)
			if count == 0 {
				break
			}
		}
		_, err = db.Query(fmt.Sprintf("insert into urls values (\"%s\", \"%s\")", short, url))
		if err != nil {
			panic(err)
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
	var err error
	dbConfig := mysql.Config{User: "root", Passwd: "anurag", Net: "tcp", Addr: "0.0.0.0:3306", DBName: "urldb"}
	db, err = sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	http.HandleFunc("/", handleIndexPage)
	http.HandleFunc("/shorten", handleShorten)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
