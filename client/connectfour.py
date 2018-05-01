import asyncio
import websockets
import json

# change to actual uri once kevin has deployed server 
uri = 'ws://35.197.29.97:8080/api/connectfour'

async def connect(uri):
    async with websockets.connect(uri) as websocket:
        # await websocket.send("{ \"hi\":\"hello\"}")
        print("connected")
        connectMsg = await websocket.recv()
        connectParsed = json.loads(connectMsg)

        if connectParsed['status'] != 'connected':
            print('Error: ' + connectParsed['message'])
            return
        
        playerNum = connectParsed['player']
        await play(websocket, playerNum)
            
async def run(websocket, player):

    while True:
        msg = await websocket.recv()
        parsed = json.loads(msg)

        if parsed['status'] != 'playing':
            if parsed['status'] == 'error':
                print('Error: ' + parsed['message'])
            
            if parsed['status'] == 'win':
                print('You won!')
            
            if parsed['status'] == 'lose':
                print('You lost :\'(')

            return
        
        board = parsed['board']

        for i in range(len(board) - 1, -1, -1):
            line = ''
            for j in range(len(board[0])):
                line += str(board[i][j]) + ' '
            print(line)

        placeColumn = play(board, player)
        print('placed in column ' + str(placeColumn) + '\n')

        await websocket.send(str(placeColumn))


"""
Complete this function, which actually plays the game.

:param board: a 6x7 integer list representing the current game board
:param player: the player number, either 1 or 2, that you are, corresponding to the numbers on board
:return: a number between 0 and 6 representing the column the program chooses to put its piece
"""
def play(board, player):
    return 6

asyncio.get_event_loop().run_until_complete(
    connect(uri))
