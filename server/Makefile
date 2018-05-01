all: update build run

update: 
	git pull

build:
	go get
	go build -o bin/server

run:
	./bin/server