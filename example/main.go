package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const closureSuffix = "/closure/goog/base.js"

var closureBaseJs = flag.String("c", "./closure-library"+closureSuffix,
	"path of closure library's base.js file")

type DataType struct {
	A int
	B string
}

var Data = DataType{1, "<b>two</b>"}

func main() {
	flag.Parse()
	fd, err := os.Open(*closureBaseJs)
	if err != nil {
		log.Fatal("cannot open " + *closureBaseJs)
	}
	fd.Close()

	t := template.Must(template.ParseFiles("index.html", "test.html"))

	closureRootPath := strings.TrimSuffix(*closureBaseJs, closureSuffix)
	closureRootURL := "/closure-library/"
	http.Handle(closureRootURL, http.StripPrefix(closureRootURL,
		http.FileServer(http.Dir(closureRootPath))))

	gotmpljsRootPath := "../js"
	gotmpljsRootURL := "/gotmpljs/"
	http.Handle(gotmpljsRootURL, http.StripPrefix(gotmpljsRootURL,
		http.FileServer(http.Dir(gotmpljsRootPath))))

	http.HandleFunc("/template.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template.js")
	})

	http.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(Data)
		w.Write(b)
	})

	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, Data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
	})

	listenAddr := ":8080"
	log.Print("listening on " + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
