install:
	go mod download

test:
	SELECTED_MODULES=${SELECTED_MODULES} poetry run python -m pytest -n auto

build:
	go build main.go

run:
	./main

build_and_run: build run