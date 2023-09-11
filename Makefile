install:
	go mod download

test:
	go test salami/compiler/lexer salami/compiler/parser salami/compiler/semantic_analyzer