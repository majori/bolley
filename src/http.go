package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Listen() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	match, _ := parseSpreadsheet(bytes.NewReader(buff))
	saveMatch(match)
	fmt.Fprint(w, "OK")
}
