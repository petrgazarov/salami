install:
	go mod download

test: # Resolves all test files ending in _test.go and runs them
	find . -name "*_test.go" -print0 | xargs -0 -n1 dirname | sort -u | xargs -L1 go test

test_coverage:
	go test ./... -coverprofile=coverage.out

build:
	go build -o salami cli/main.go