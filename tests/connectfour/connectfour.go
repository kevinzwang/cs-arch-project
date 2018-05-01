package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var addr = "localhost:8080"
var upgrader = websocket.Upgrader{}
var players = make([]*websocket.Conn, 2)

func main() {
	http.HandleFunc("/api/connectfour", connect)
	http.HandleFunc("/home", home)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hi")
}

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	if players[0] == nil {
		fmt.Println("Player 1 connected")
		c.WriteJSON(map[string]interface{}{
			"status": "connected",
			"player": 1,
		})
		players[0] = c
	} else if players[1] == nil {
		fmt.Println("Player 2 connected")
		c.WriteJSON(map[string]interface{}{
			"status": "connected",
			"player": 2,
		})
		players[1] = c
		play()
	} else {
		c.WriteJSON(map[string]interface{}{
			"message": "Already two players",
			"status":  "error",
		})
		c.Close()
	}
}

func play() {
	fmt.Println("Starting game")
	for _, p := range players {
		defer p.Close()
	}

	board := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}

	for turn := 0; ; turn = (turn + 1) % 2 {
		fmt.Println("Player " + strconv.Itoa(turn+1) + " turn")
		curr := players[turn]
		other := players[(turn+1)%2]

		curr.WriteJSON(map[string]interface{}{
			"board":   board,
			"message": "Your turn!",
			"status":  "playing",
		})

		_, message, err := curr.ReadMessage()
		if err != nil {
			break
		}

		column, err := strconv.Atoi(string(message))
		if err != nil {
			curr.WriteJSON(map[string]interface{}{
				"message": "Not a number",
				"status":  "error",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "error",
			})
			return
		}
		if column < 0 || column > 6 {
			curr.WriteJSON(map[string]interface{}{
				"message": "Number must be in range 0 to 6.",
				"status":  "error",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "error",
			})
			return
		}

		row := -1
		for i := 0; i < len(board); i++ {
			if board[i][column] == 0 {
				board[i][column] = turn + 1
				row = i
				break
			}
		}

		if row == -1 {
			curr.WriteJSON(map[string]interface{}{
				"message": "Column already full",
				"status":  "error",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "error",
			})
			return
		}

		if checkWin(board, turn+1, row, column) {
			curr.WriteJSON(map[string]interface{}{
				"status": "win",
				"board":  board,
			})
			other.WriteJSON(map[string]interface{}{
				"status": "lose",
				"board":  board,
			})
			return
		}
	}
}

var directions = [][]int{
	{0, 1},
	{1, 0},
	{1, 1},
	{1, -1},
}

func checkWin(board [][]int, player int, row int, column int) bool {
	for _, d := range directions {
		if checkDirection(board, player, row, column, d[0], d[1])+checkDirection(board, player, row, column, -1*d[0], -1*d[1]) >= 3 {
			return true
		}
	}
	return false
}

func checkDirection(board [][]int, player int, row int, column int, dr int, dc int) (count int) {
	for count = -1; row >= 0 && row < len(board) && column >= 0 && column < len(board[0]) && count < 3; row, column, count = row+dr, column+dc, count+1 {
		if board[row][column] != player {
			break
		}
	}
	return
}
