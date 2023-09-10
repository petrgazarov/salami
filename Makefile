install:
	go mod download

test:
	go test salami/compiler/lexer salami/compiler/parser salami/compiler/semantic_analyzer

build:
	go build main.go

run:
	./main

build_and_run: build run