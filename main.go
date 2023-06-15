package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	_ "modernc.org/sqlite"
)

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/static/style.css", style)
	http.HandleFunc("/create", create)
	log.Println("Serving on http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func root(w http.ResponseWriter, req *http.Request) {
    // isn't plain root
    if !strings.HasSuffix(req.URL.Path, "/") {
        log.Println("url forwarding")
        id := req.URL.Path[1:]
        db, err := sql.Open("sqlite", "url.db")
        if err != nil {
            log.Fatalf("%v", err)
        }
        var original string
        err = db.QueryRow("SELECT original FROM shortener WHERE id = $1", id).Scan(&original);
        if err != nil {
            log.Printf("%v", err)
            w.WriteHeader(404)
            return
        }

        w.Header().Add("Location", original)
        w.WriteHeader(302)

        return
    }
    log.Println("serving root")
	bytes, err := os.ReadFile("static/index.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf8")
	w.Write(bytes)
}

func style(w http.ResponseWriter, req *http.Request) {
	bytes, err := os.ReadFile("static/style.css")
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.Header().Add("Content-Type", "text/css; charset=utf8")
	w.Write(bytes)
}

func create(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	url := req.Form.Get("url")
	if url == "" {
        w.Header().Add("Location", "http://localhost:8000/")
        w.WriteHeader(302)
		return
	}
	log.Printf("shortening %s", url)

	db, err := sql.Open("sqlite", "url.db")
	if err != nil {
		log.Fatalf("%v", err)
	}

	var id int
	row := db.QueryRow("INSERT INTO shortener (original) VALUES ($1) RETURNING id", url)
	err = row.Scan(&id)
	if err != nil {
		log.Fatalf("%v", err)
	}

	bytes, err := os.ReadFile("static/create.html")
	str := string(bytes)
    str = strings.ReplaceAll(str, "{url}", fmt.Sprintf("http://localhost:8000/%d", id))
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf8")
	io.WriteString(w, str)
}
