package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/static/style.css", style)
	http.HandleFunc("/static/index.js", script)
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
		var expires string
		err = db.QueryRow("SELECT original, expires FROM shortener WHERE hash = $1", id).Scan(&original, &expires)
		if err != nil {
			log.Printf("selecting url: %v", err)
			w.WriteHeader(404)
			return
		}

		t, err := time.Parse("2006-01-02T15:04", expires)
		if err != nil {
			log.Printf("date invalid: %v", err)
			return
		}
		if time.Now().After(t) {
			w.WriteHeader(410)
			return
		}

		w.Header().Add("Location", original)
		w.WriteHeader(302)

		return
	}
	log.Println("serving root")
	bytes, err := os.ReadFile("static/index.html")
	if err != nil {
		log.Fatalf("could not open index.html: %v", err)
	}
	w.Header().Add("Content-Type", "text/html; charset=utf8")
	w.Write(bytes)
}

func style(w http.ResponseWriter, req *http.Request) {
	bytes, err := os.ReadFile("static/style.css")
	if err != nil {
		log.Fatalf("could not open style.css: %v", err)
	}
	w.Header().Add("Content-Type", "text/css; charset=utf8")
	w.Write(bytes)
}

func script(w http.ResponseWriter, req *http.Request) {
	bytes, err := os.ReadFile("static/index.js")
	if err != nil {
		log.Fatalf("could not open index.js: %v", err)
	}
	w.Header().Add("Content-Type", "application/javascript; charset=utf8")
	w.Write(bytes)
}

func create(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	url := req.Form.Get("url")
	expires := req.Form.Get("expire")
	if url == "" || expires == "" {
		w.Header().Add("Location", "http://localhost:8000/")
		w.WriteHeader(302)
		return
	}
	log.Printf("shortening %s", url)

	db, err := sql.Open("sqlite", "url.db")
	if err != nil {
		log.Fatalf("err opening db: %v", err)
	}

	hasher := md5.New()
	_, err = io.WriteString(hasher, url)
	if err != nil {
		log.Fatalf("could not write to md5 hash")
	}
	md5 := hasher.Sum(nil)
	hash := base64.StdEncoding.EncodeToString(md5)[:7]
	_, err = db.Exec("INSERT INTO shortener (original, hash, expires) VALUES ($1, $2, $3)", url, hash, expires)
	if err != nil {
		log.Fatalf("err inserting: %v", err)
	}

	bytes, err := os.ReadFile("static/create.html")
	if err != nil {
		log.Fatalf("err reading create.html: %v", err)
	}
	str := string(bytes)
	str = strings.ReplaceAll(str, "{url}", fmt.Sprintf("http://localhost:8000/%s", hash))
	w.Header().Add("Content-Type", "text/html; charset=utf8")
	io.WriteString(w, str)
}
