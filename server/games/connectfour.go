package games

import (
	"time"

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

func (game *ConnectFour) play(p1 *websocket.Conn, p2 *websocket.Conn) *websocket.Conn {
	players := [2]*websocket.Conn{p1, p2}
	for i, p := range players {
		p.WriteJSON(map[string]interface{}{
			"status": "start",
			"player": i + 1,
		})
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

	for turn, count := 0, 0; count < 42; turn, count = (turn+1)%2, count+1 {
		curr := players[turn]
		other := players[(turn+1)%2]

		curr.WriteJSON(map[string]interface{}{
			"board":   board,
			"message": "Your turn!",
			"status":  "playing",
		})

		var msg map[string]interface{}
		resp := make(chan map[string]interface{})
		go chanReadJSON(curr, resp)

		select {
		case msg = <-resp:
		case <-time.After(15 * time.Second):
			curr.WriteJSON(map[string]interface{}{
				"message": "took too long to respond",
				"status":  "failure",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent took too long to respond",
				"status":  "failure",
			})
			break
		}

		foo, ok := msg["move"].(float64)
		column := int(foo)
		if !ok {
			curr.WriteJSON(map[string]interface{}{
				"message": "Not a number",
				"status":  "failure",
			})
			other.WriteJSON(map[string]interface{}{
				"message": "Opponent dun a goof",
				"status":  "failure",
			})
			return other
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
			return other
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
			return other
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
			return curr
		}
	}

	for _, p := range players {
		p.WriteJSON(map[string]interface{}{
			"status": "tie",
			"board":  board,
		})
	}
	return nil
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
