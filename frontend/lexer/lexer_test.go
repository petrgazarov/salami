package lexer_test

import (
	"os"
	"path/filepath"
	"salami/frontend/lexer"
	"salami/frontend/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %s", err)
	}
	fixturePath := filepath.Join(cwd, "testdata", "source.sami")

	source, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("Failed to read fixture: %s", err)
	}

	lexer := lexer.NewLexer(fixturePath, string(source))
	tokens, err := lexer.Run()
	if err != nil {
		t.Fatalf("Failed to run lexer: %s", err)
	}
	expectedTokens := getExpectedTokens()

	assert.Equal(t, len(expectedTokens), len(tokens), "Unexpected number of tokens")
	for i, token := range tokens {
		require.Equal(t, expectedTokens[i], token, "Unexpected token at %d", i)
	}
}

func getExpectedTokens() []*types.Token {
	return []*types.Token{
		{Type: types.DecoratorName, Value: "@resource", Line: 1, Column: 1},
		{Type: types.DecoratorArg, Value: "aws.cloudwatch.LogGroup", Line: 1, Column: 11},
		{Type: types.DecoratorArg, Value: "ServerLogGroup", Line: 1, Column: 36},
		{Type: types.Newline, Value: "", Line: 1, Column: 51},
		{Type: types.NaturalLanguage, Value: "Name: server-log-group", Line: 2, Column: 1},
		{Type: types.Newline, Value: "", Line: 2, Column: 23},
		{Type: types.Newline, Value: "", Line: 3, Column: 1},
		{Type: types.DecoratorName, Value: "@resource", Line: 4, Column: 1},
		{Type: types.DecoratorArg, Value: "aws.ecr.Repository", Line: 4, Column: 11},
		{Type: types.DecoratorArg, Value: "ServerRepository", Line: 4, Column: 31},
		{Type: types.Newline, Value: "", Line: 4, Column: 48},
		{Type: types.NaturalLanguage, Value: "Name: server", Line: 5, Column: 1},
		{Type: types.Newline, Value: "", Line: 5, Column: 13},
		{Type: types.NaturalLanguage, Value: "Has mutable image tags", Line: 6, Column: 1},
		{Type: types.Newline, Value: "", Line: 6, Column: 23},
		{Type: types.Newline, Value: "", Line: 7, Column: 1},
		{Type: types.DecoratorName, Value: "@resource", Line: 8, Column: 1},
		{Type: types.DecoratorArg, Value: "aws.ecr.LifecyclePolicy", Line: 8, Column: 11},
		{Type: types.DecoratorArg, Value: "ServerRepoLifecyclePolicy", Line: 8, Column: 36},
		{Type: types.Newline, Value: "", Line: 8, Column: 62},
		{Type: types.NaturalLanguage, Value: "In $ServerRepository", Line: 9, Column: 1},
		{Type: types.Newline, Value: "", Line: 9, Column: 21},
		{
			Type:   types.NaturalLanguage,
			Value:  "Policy retains only the last 10 untagged images in the repository. Images beyond this count will expire.",
			Line:   10,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 10, Column: 105},
		{Type: types.Newline, Value: "", Line: 11, Column: 1},
		{Type: types.DecoratorName, Value: "@resource", Line: 12, Column: 1},
		{Type: types.DecoratorArg, Value: "aws.ecs.Service", Line: 12, Column: 11},
		{Type: types.DecoratorArg, Value: "ServerEcsService", Line: 12, Column: 28},
		{Type: types.Newline, Value: "", Line: 12, Column: 45},
		{
			Type:   types.NaturalLanguage,
			Value:  "In $EcsCluster, using $ServerTaskDefinition",
			Line:   13,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 13, Column: 44},
		{Type: types.NaturalLanguage, Value: "Name: server", Line: 14, Column: 1},
		{Type: types.Newline, Value: "", Line: 14, Column: 13},
		{Type: types.NaturalLanguage, Value: "Desired count: 1", Line: 15, Column: 1},
		{Type: types.Newline, Value: "", Line: 15, Column: 17},
		{Type: types.NaturalLanguage, Value: "Launch type: FARGATE", Line: 16, Column: 1},
		{Type: types.Newline, Value: "", Line: 16, Column: 21},
		{Type: types.NaturalLanguage, Value: "---", Line: 17, Column: 1},
		{Type: types.Newline, Value: "", Line: 17, Column: 4},
		{Type: types.NaturalLanguage, Value: "Network configuration:", Line: 18, Column: 1},
		{Type: types.Newline, Value: "", Line: 18, Column: 23},
		{Type: types.NaturalLanguage, Value: "  Assigned public IP", Line: 19, Column: 1},
		{Type: types.Newline, Value: "", Line: 19, Column: 21},
		{Type: types.NaturalLanguage, Value: "  Subnets: $PublicSubnetA and $PublicSubnetB", Line: 20, Column: 1},
		{Type: types.Newline, Value: "", Line: 20, Column: 45},
		{Type: types.NaturalLanguage, Value: "  Security group: $ServerEcsSecurityGroup", Line: 21, Column: 1},
		{Type: types.Newline, Value: "", Line: 21, Column: 42},
		{Type: types.NaturalLanguage, Value: "Load balancers:", Line: 22, Column: 1},
		{Type: types.Newline, Value: "", Line: 22, Column: 16},
		{Type: types.NaturalLanguage, Value: "  Target group: $ServerTargetGroup", Line: 23, Column: 1},
		{Type: types.Newline, Value: "", Line: 23, Column: 35},
		{Type: types.NaturalLanguage, Value: "  Container name: {server_container_name}", Line: 24, Column: 1},
		{Type: types.Newline, Value: "", Line: 24, Column: 42},
		{Type: types.NaturalLanguage, Value: "  Port: {container_port}", Line: 25, Column: 1},
		{Type: types.Newline, Value: "", Line: 25, Column: 25},
		{Type: types.NaturalLanguage, Value: "Deployment:", Line: 26, Column: 1},
		{Type: types.Newline, Value: "", Line: 26, Column: 12},
		{
			Type:   types.NaturalLanguage,
			Value:  "  ECS type deployment controller",
			Line:   27,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 27, Column: 33},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Deployment circuit breaker: enabled with rollback",
			Line:   28,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 28, Column: 52},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Wait for steady state: True",
			Line:   29,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 29, Column: 30},
		{Type: types.Newline, Value: "", Line: 30, Column: 1},
		{Type: types.DecoratorName, Value: "@variable", Line: 31, Column: 1},
		{Type: types.DecoratorArg, Value: "server_container_name", Line: 31, Column: 11},
		{Type: types.DecoratorArg, Value: "string", Line: 31, Column: 34},
		{Type: types.DecoratorArg, Value: "server-container", Line: 31, Column: 42},
		{Type: types.Newline, Value: "", Line: 31, Column: 59},
		{Type: types.NaturalLanguage, Value: "Description: Server container name", Line: 32, Column: 1},
		{Type: types.EOF, Value: "", Line: 32, Column: 35},
	}
}
