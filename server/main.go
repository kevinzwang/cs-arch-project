package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/hi", hiHandler)

	http.Handle("/", r)
	appengine.Main()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hi")
}

func hiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!")
}
