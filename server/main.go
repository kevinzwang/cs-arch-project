package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"./games"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var tokens = make(map[int64]bool)

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("website/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.HandleFunc("/api/{game}", GameHandler)
	r.HandleFunc("/view/{game}/{token}", ViewHandler)
	r.HandleFunc("/new-token", TokenGenerateHandler)
	r.HandleFunc("/list-tokens", ListTokenHandler)
	r.HandleFunc("/docs/{game}", DocsHandler)
	r.HandleFunc("/", HomeHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HomeHandler handles the home page for the site
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("website/templates/index.html"))
	tmpl.Execute(w, nil)
}

// ViewHandler allows you to view the game
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Game: "+mux.Vars(r)["game"])
	fmt.Fprintln(w, "Token: "+mux.Vars(r)["token"])
}

func DocsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("website/templates/docs/" + mux.Vars(r)["game"] + ".html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tmpl.Execute(w, nil)
}

// GameHandler handles API calls for games
func GameHandler(w http.ResponseWriter, r *http.Request) {
	gameName := mux.Vars(r)["game"]

	var gameStruct games.Game
	switch gameName {
	case "connectfour":
		gameStruct = &games.ConnectFour{}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error: ", err)
		return
	}

	token, ok := authenticate(conn)
	if !ok {
		return
	}
	defer func() { tokens[token] = false }()

	gameStruct.Execute(conn)
}

func authenticate(conn *websocket.Conn) (int64, bool) {
	var msg map[string]interface{}
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Print("read json error: ", err)
		conn.WriteJSON(map[string]interface{}{
			"status":  "failure",
			"message": "Couldn't read JSON",
		})
		return 0, false
	}

	token, ok := msg["token"]
	if !ok {
		conn.WriteJSON(map[string]interface{}{
			"status":  "failure",
			"message": "Couldn't get token from JSON",
		})
		return 0, false
	}

	fmtToken := int64(token.(float64))

	playing, ok := tokens[fmtToken]

	if !ok {
		conn.WriteJSON(map[string]interface{}{
			"status":  "failure",
			"message": "No such token",
		})
		return 0, false
	}

	if playing {
		conn.WriteJSON(map[string]interface{}{
			"status":  "failure",
			"message": "Already an instance playing",
		})
		return 0, false
	}

	dt := time.Now().Unix() - fmtToken
	if dt > 172800 {
		delete(tokens, fmtToken)
		conn.WriteJSON(map[string]interface{}{
			"status":  "failure",
			"message": "Token has expired",
		})
		return 0, false
	}

	time.Sleep(2 * time.Second)
	tokens[fmtToken] = true
	conn.WriteJSON(map[string]interface{}{
		"status": "connected",
	})
	return fmtToken, true
}

// TokenGenerateHandler generates tokens to use
func TokenGenerateHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Unix()
	_, err := fmt.Fprint(w, t)
	if err != nil {
		log.Print("http print error: ", err)
		return
	}
	tokens[t] = false
}

// ListTokenHandler lists the tokens that are currently available
func ListTokenHandler(w http.ResponseWriter, r *http.Request) {
	for t := range tokens {
		fmt.Fprintln(w, t)
	}
}
