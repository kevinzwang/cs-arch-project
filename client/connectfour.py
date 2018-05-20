import asyncio
import websockets
import json

uri = 'ws://localhost:8080/api/connectfour'
# token = 1526776854

async def connect(uri):
    async with websockets.connect(uri) as websocket:
        token = input("token: ")
        await websocket.send(json.dumps({"token": int(token)}))
        connectMsg = await websocket.recv()
        connectParsed = json.loads(connectMsg)

        if connectParsed['status'] != 'connected':
            print('Failure to connect: ' + connectParsed['message'])
            return
        
        print('connected!')
        
        startMsg = await websocket.recv()
        startParsed = json.loads(startMsg)
        if startParsed['status'] != 'start':
            print('Failure to start: ' + startParsed['message'])
            return
        
        player = int(startParsed['player'])
        print('player: ' + str(player) + '\n')

        while True:
            turnMsg = await websocket.recv()
            turnParsed = json.loads(turnMsg)
            
            if turnParsed['status'] != 'playing':
                if turnParsed['status'] == 'win':
                    print('you won!')
                elif turnParsed['status'] == 'lose':
                    print('you lost.')
                elif turnParsed['status'] == 'failure':
                    print('failure: ' + turnParsed['message'])
                else:
                    print('unknown problem: ' + turnParsed['status'] + ' - ' + turnParsed['message'])
                return

            printboard(turnParsed['board'])
            
            move = play(player, turnParsed['board'])
            await websocket.send(json.dumps({'move': move}))

def printboard(board):
    for row in reversed(board):
        r = ''
        for column in row:
            r += str(column)
        print(r)
    print()

"""
Complete this function, which actually plays the game.

:param board: a 6x7 integer list representing the current game board
:param player: the player number, either 1 or 2, that you are, corresponding to the numbers on board
:return: a number between 0 and 6 representing the column the program chooses to put its piece
"""
def play(player, board):
    return 6

asyncio.get_event_loop().run_until_complete(
    connect(uri))
