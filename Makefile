install:
	go mod download

test:
	go test salami/compiler/lexer

build:
	go build main.go

run:
	./main

build_and_run: build run