// Example code for the GoTmplJs compiler.
// Copyright (C) 2014  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
	B []string
}

var Data = DataType{
	A: 1,
	B: []string{"one", "<b>two</b>", "three"},
}

func main() {
	flag.Parse()
	fd, err := os.Open(*closureBaseJs)
	if err != nil {
		log.Fatal("cannot open " + *closureBaseJs)
	}
	fd.Close()

	t := template.Must(template.ParseFiles("index.html", "example.html"))

	closureRootPath := strings.TrimSuffix(*closureBaseJs, closureSuffix)
	closureRootURL := "/closure-library/"
	http.Handle(closureRootURL, http.StripPrefix(closureRootURL,
		http.FileServer(http.Dir(closureRootPath))))
	http.HandleFunc("/gotmpl.js",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../gotmpl.js")
		})
	http.HandleFunc("/example.js",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "example.js")
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
