package games

import (
	"strconv"

	"github.com/gorilla/websocket"
)

// ConnectFour is the handler for the connect four game
type ConnectFour struct{}

// Execute puts players on queue or to face off
func (game *ConnectFour) Execute(conn *websocket.Conn) {
	if queuing["connectfour"] == nil {
		queuing["connectfour"] = conn
	} else {
		p1 := queuing["connectfour"]
		queuing["connectfour"] = nil
		game.play(p1, conn)
	}
}

// func chanReadMessage(conn *websocket.Conn, )

func (game *ConnectFour) play(p1 *websocket.Conn, p2 *websocket.Conn) int {
	players := [2]*websocket.Conn{p1, p2}
	for i, p := range players {
		p.WriteJSON(map[string]interface{}{
			"status": "start",
			"player": i + 1,
		})
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
				"status":  "failure",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "failure",
			})
			return (turn+1)%2 + 1
		}
		if column < 0 || column > 6 {
			curr.WriteJSON(map[string]interface{}{
				"message": "Number must be in range 0 to 6.",
				"status":  "failure",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "failure",
			})
			return (turn+1)%2 + 1
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
				"status":  "failure",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "failure",
			})
			return (turn+1)%2 + 1
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
			return (turn + 1) % 2
		}
	}
	return 0
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
