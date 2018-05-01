package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.HandleFunc("/api/connectfour", connectFourHandler)
	r.HandleFunc("/", homeHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
	// appengine.Main()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func connectFourHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}

	defer conn.Close()

	ping := new(map[interface{}]interface{})
	err = conn.ReadJSON(ping)
	if err != nil {
		log.Print("read json error:", err)
		return
	}

	pingAsJSON, _ := json.Marshal(ping)
	fmt.Println(pingAsJSON)
}
