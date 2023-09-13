install:
	go mod download

test:
	go test salami/frontend/lexer salami/frontend/parser salami/frontend/semantic_analyzer