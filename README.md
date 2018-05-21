# CS Arch Final Project - Bot Game
by Kevin Wang

## Intro

This repo contains the server and client implementations of a framework that I made to allow people to easily create "bots" for games such as connect 4 and tetris. 

This framework uses Websockets with JSON to transmit information between the client and server. 

## Limitations

Due to the time constraints, only the Connect 4 game is completed, and authentication is with generated tokens instead of creating a login/signup system since that is not integral to the concept and requires knowledge about securely transmitting, storing, and verifying usernames and passwords which I do not want to be trusted upon in just one quarter of a school year. 

## Getting Started

The server is hosted on http://35.197.29.97:8080/. A full implementation of the provided API is in `client/connectfour.py`, which requires Python 3.5+ to run and also the package `websockets` which can be installed with `python3 -m pip install websockets`. All you need to do to play is implement `def play()`, grab a token from the website and set the variable `token` to it, and run it.


If you would like to host the server yourself, you first need the latest version of Golang (can be found [here](https://golang.org/)) and the command line tool `make`. Navigate to the `server` folder and run 

## Contents

- _server_- the server-side code, written in Golang
    - **main.go** - the main server file. contains the http server and url handlers and connects the different parts of the server application
    - **Makefile** - you can use this to automatically download the latest version, build, and/or run the game
    - _games_ - the games package for the API
        - **connectfour.go** - contains all the connect four-specific components
        - **games.go** - contains the games
    - _website_ - the HTML, CSS, and image files used for the website
        - _assets_ - the static contents, such as CSS and images
        - _templates_ - the HTML templates
            - **index.html** - the home page for the website
            - _docs_ - API documentation pages
                - **connectfour.html** - documentation for the connect four game
- _client_ - where the example client API implementations can be found
    - **connectfour<nolink>.py** - a working connect four implementation
- _tests_ - this folder is just for when I was testing the different libraries and frameworks pls no look xd

