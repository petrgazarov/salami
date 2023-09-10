package lexer_test

import (
	"os"
	"path/filepath"
	"salami/compiler/lexer"
	"salami/compiler/types"
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
		{Type: types.FieldName, Value: "Resource type", Line: 1, Column: 1},
		{Type: types.FieldValue, Value: "cloudwatch.LogGroup", Line: 1, Column: 16},
		{Type: types.Newline, Value: "", Line: 1, Column: 35},
		{Type: types.FieldName, Value: "Logical name", Line: 2, Column: 1},
		{Type: types.FieldValue, Value: "CumuliServerLogGroup", Line: 2, Column: 15},
		{Type: types.Newline, Value: "", Line: 2, Column: 35},
		{Type: types.FieldName, Value: "Name", Line: 3, Column: 1},
		{Type: types.FieldValue, Value: "cumuli-server-log-group", Line: 3, Column: 7},
		{Type: types.Newline, Value: "", Line: 3, Column: 30},
		{Type: types.Newline, Value: "", Line: 4, Column: 1},
		{Type: types.DecoratorName, Value: "@exports", Line: 5, Column: 1},
		{
			Type:   types.DecoratorArg,
			Value:  "name: cumuli-server-ecr-repository-name",
			Line:   5,
			Column: 10,
		},
		{Type: types.Newline, Value: "", Line: 5, Column: 50},
		{Type: types.FieldName, Value: "Resource type", Line: 6, Column: 1},
		{Type: types.FieldValue, Value: "ecr.Repository", Line: 6, Column: 16},
		{Type: types.Newline, Value: "", Line: 6, Column: 30},
		{Type: types.FieldName, Value: "Logical name", Line: 7, Column: 1},
		{Type: types.FieldValue, Value: "CumuliServerRepository", Line: 7, Column: 15},
		{Type: types.Newline, Value: "", Line: 7, Column: 37},
		{Type: types.FieldName, Value: "Name", Line: 8, Column: 1},
		{Type: types.FieldValue, Value: "cumuli-server", Line: 8, Column: 7},
		{Type: types.Newline, Value: "", Line: 8, Column: 20},
		{Type: types.NaturalLanguage, Value: "Has mutable image tags.", Line: 9, Column: 1},
		{Type: types.Newline, Value: "", Line: 9, Column: 24},
		{Type: types.Newline, Value: "", Line: 10, Column: 1},
		{Type: types.DecoratorName, Value: "@uses", Line: 11, Column: 1},
		{Type: types.DecoratorArg, Value: "CumuliServerRepository", Line: 11, Column: 7},
		{Type: types.Newline, Value: "", Line: 11, Column: 30},
		{Type: types.FieldName, Value: "Resource type", Line: 12, Column: 1},
		{Type: types.FieldValue, Value: "ecr.LifecyclePolicy", Line: 12, Column: 16},
		{Type: types.Newline, Value: "", Line: 12, Column: 35},
		{Type: types.FieldName, Value: "Logical name", Line: 13, Column: 1},
		{Type: types.FieldValue, Value: "CumuliServerRepoLifecyclePolicy", Line: 13, Column: 15},
		{Type: types.Newline, Value: "", Line: 13, Column: 46},
		{
			Type: types.NaturalLanguage,
			Value: "Policy: A JSON policy with a rule that retains only the last 10 untagged images in the repository. " +
				"Images beyond this count will expire.",
			Line:   14,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 14, Column: 137},
		{Type: types.Newline, Value: "", Line: 15, Column: 1},
		{Type: types.DecoratorName, Value: "@exports", Line: 16, Column: 1},
		{Type: types.DecoratorArg, Value: "name: exported-name", Line: 16, Column: 10},
		{Type: types.Newline, Value: "", Line: 16, Column: 30},
		{Type: types.DecoratorName, Value: "@uses", Line: 17, Column: 1},
		{Type: types.DecoratorArg, Value: "CumuliEcsCluster", Line: 17, Column: 7},
		{Type: types.DecoratorArg, Value: "CumuliServerTaskDefinition", Line: 17, Column: 25},
		{Type: types.DecoratorArg, Value: "PublicSubnetA", Line: 17, Column: 53},
		{Type: types.DecoratorArg, Value: "PublicSubnetB", Line: 17, Column: 68},
		{Type: types.DecoratorArg, Value: "CumuliServerEcsSecurityGroup", Line: 17, Column: 83},
		{Type: types.DecoratorArg, Value: "CumuliServerTargetGroup", Line: 17, Column: 113},
		{Type: types.Newline, Value: "", Line: 17, Column: 137},
		{Type: types.FieldName, Value: "Resource type", Line: 18, Column: 1},
		{Type: types.FieldValue, Value: "aws.ecs.Service", Line: 18, Column: 16},
		{Type: types.Newline, Value: "", Line: 18, Column: 31},
		{Type: types.FieldName, Value: "Logical name", Line: 19, Column: 1},
		{Type: types.FieldValue, Value: "CumuliServerEcsService", Line: 19, Column: 15},
		{Type: types.Newline, Value: "", Line: 19, Column: 37},
		{Type: types.FieldName, Value: "Name", Line: 20, Column: 1},
		{Type: types.FieldValue, Value: "cumuli-server", Line: 20, Column: 7},
		{Type: types.Newline, Value: "", Line: 20, Column: 20},
		{Type: types.NaturalLanguage, Value: "Desired count: 1", Line: 21, Column: 1},
		{Type: types.Newline, Value: "", Line: 21, Column: 17},
		{Type: types.NaturalLanguage, Value: "Launch type: FARGATE", Line: 22, Column: 1},
		{Type: types.Newline, Value: "", Line: 22, Column: 21},
		{Type: types.NaturalLanguage, Value: "---", Line: 23, Column: 1},
		{Type: types.Newline, Value: "", Line: 23, Column: 4},
		{Type: types.NaturalLanguage, Value: "Network configuration:", Line: 24, Column: 1},
		{Type: types.Newline, Value: "", Line: 24, Column: 23},
		{Type: types.NaturalLanguage, Value: "  Assigned public IP.", Line: 25, Column: 1},
		{Type: types.Newline, Value: "", Line: 25, Column: 22},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Subnets: PublicSubnetA and PublicSubnetB.",
			Line:   26,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 26, Column: 44},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Security group: CumuliServerEcsSecurityGroup.",
			Line:   27,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 27, Column: 48},
		{Type: types.NaturalLanguage, Value: "Load balancers:", Line: 28, Column: 1},
		{Type: types.Newline, Value: "", Line: 28, Column: 16},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Target group: CumuliServerTargetGroup.",
			Line:   29,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 29, Column: 41},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Container name: {server_container_name}.",
			Line:   30,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 30, Column: 43},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Port: {container_port}.",
			Line:   31,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 31, Column: 26},
		{Type: types.NaturalLanguage, Value: "Deployment:", Line: 32, Column: 1},
		{Type: types.Newline, Value: "", Line: 32, Column: 12},
		{
			Type:   types.NaturalLanguage,
			Value:  "  ECS type deployment controller.",
			Line:   33,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 33, Column: 34},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Deployment circuit breaker: enabled with rollback.",
			Line:   34,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 34, Column: 53},
		{
			Type:   types.NaturalLanguage,
			Value:  "  Wait for steady state: True",
			Line:   35,
			Column: 1,
		},
		{Type: types.Newline, Value: "", Line: 35, Column: 30},
		{Type: types.Newline, Value: "", Line: 36, Column: 1},
		{Type: types.DecoratorName, Value: "@variable", Line: 37, Column: 1},
		{Type: types.DecoratorArg, Value: "string", Line: 37, Column: 11},
		{Type: types.Newline, Value: "", Line: 37, Column: 18},
		{Type: types.FieldName, Value: "Description", Line: 38, Column: 1},
		{Type: types.FieldValue, Value: "Server container name", Line: 38, Column: 14},
		{Type: types.Newline, Value: "", Line: 38, Column: 35},
		{Type: types.FieldName, Value: "Name", Line: 39, Column: 1},
		{Type: types.FieldValue, Value: "server_container_name", Line: 39, Column: 7},
		{Type: types.Newline, Value: "", Line: 39, Column: 28},
		{Type: types.FieldName, Value: "Value", Line: 40, Column: 1},
		{Type: types.FieldValue, Value: "cumuli-server-container", Line: 40, Column: 8},
		{Type: types.EOF, Value: "", Line: 40, Column: 31},
	}
}
