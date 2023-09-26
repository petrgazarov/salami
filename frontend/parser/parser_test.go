package parser_test

import (
	commonTypes "salami/common/types"
	"salami/frontend/parser"
	frontendTypes "salami/frontend/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	inputTokens := getInput()
	dummyFilePath := "dummy/file/path"
	parser := parser.NewParser(inputTokens, dummyFilePath)
	resources, variables, err := parser.Parse()
	if err != nil {
		t.Fatalf(err.Error())
	}
	expectedResources := getExpectedResources()
	expectedVariables := getExpectedVariables(t)
	assert.Equal(t, len(expectedResources), len(resources), "Unexpected number of resources")
	assert.Equal(t, len(expectedVariables), len(variables), "Unexpected number of variables")
	for i, resource := range resources {
		require.Equal(t, expectedResources[i], resource, "Unexpected resource at %d", i)
	}
}

func getExpectedResources() []*commonTypes.ParsedResource {
	return []*commonTypes.ParsedResource{
		{
			ResourceType:        "aws.cloudwatch.LogGroup",
			LogicalName:         "ServerLogGroup",
			NaturalLanguage:     "Name: server-log-group",
			ReferencedResources: []commonTypes.LogicalName{},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
			SourceFileLine:      1,
		},
		{
			ResourceType:        "aws.ecr.Repository",
			LogicalName:         "ServerRepository",
			NaturalLanguage:     "Name: server\nHas mutable image tags",
			ReferencedResources: []commonTypes.LogicalName{},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
			SourceFileLine:      4,
		},
		{
			ResourceType:        "aws.ecr.LifecyclePolicy",
			LogicalName:         "ServerRepoLifecyclePolicy",
			NaturalLanguage:     "In $ServerRepository\nPolicy retains only the last 10 untagged images in the repository. Images beyond this count will expire.",
			ReferencedResources: []commonTypes.LogicalName{"ServerRepository"},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
			SourceFileLine:      8,
		},
		{
			ResourceType: "aws.ecs.Service",
			LogicalName:  "ServerEcsService",
			NaturalLanguage: `In $EcsCluster, using $ServerTaskDefinition
Name: server
Desired count: 1
Launch type: FARGATE
---
Network configuration:
  Assigned public IP
  Subnets: $PublicSubnetA and $PublicSubnetB
  Security group: $ServerEcsSecurityGroup
Load balancers:
  Target group: $ServerTargetGroup
  Container name: {server_container_name}
  Port: {container_port}
Deployment:
  ECS type deployment controller
  Deployment circuit breaker: enabled with rollback
  Wait for steady state: True`,
			ReferencedResources: []commonTypes.LogicalName{
				"EcsCluster",
				"ServerTaskDefinition",
				"PublicSubnetA",
				"PublicSubnetB",
				"ServerEcsSecurityGroup",
				"ServerTargetGroup",
			},
			ReferencedVariables: []string{"server_container_name", "container_port"},
			SourceFilePath:      "dummy/file/path",
			SourceFileLine:      12,
		},
	}
}

func getExpectedVariables(t *testing.T) []*commonTypes.ParsedVariable {
	variableType := commonTypes.VariableType("string")
	return []*commonTypes.ParsedVariable{
		{
			Name:            "server_container_name",
			NaturalLanguage: "Description: Server container name",
			Default:         "server-container",
			Type:            variableType,
			SourceFilePath:  "dummy/file/path",
			SourceFileLine:  24,
		},
	}
}

func getInput() []*frontendTypes.Token {
	return []*frontendTypes.Token{
		{Type: frontendTypes.DecoratorName, Value: "@resource", Line: 1, Column: 1},
		{Type: frontendTypes.DecoratorArg, Value: "aws.cloudwatch.LogGroup", Line: 1, Column: 11},
		{Type: frontendTypes.DecoratorArg, Value: "ServerLogGroup", Line: 1, Column: 36},
		{Type: frontendTypes.Newline, Value: "", Line: 1, Column: 51},
		{Type: frontendTypes.NaturalLanguage, Value: "Name: server-log-group", Line: 2, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 2, Column: 23},
		{Type: frontendTypes.Newline, Value: "", Line: 3, Column: 1},
		{Type: frontendTypes.DecoratorName, Value: "@resource", Line: 4, Column: 1},
		{Type: frontendTypes.DecoratorArg, Value: "aws.ecr.Repository", Line: 4, Column: 11},
		{Type: frontendTypes.DecoratorArg, Value: "ServerRepository", Line: 4, Column: 31},
		{Type: frontendTypes.Newline, Value: "", Line: 4, Column: 48},
		{Type: frontendTypes.NaturalLanguage, Value: "Name: server", Line: 5, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 5, Column: 13},
		{Type: frontendTypes.NaturalLanguage, Value: "Has mutable image tags", Line: 6, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 6, Column: 23},
		{Type: frontendTypes.Newline, Value: "", Line: 7, Column: 1},
		{Type: frontendTypes.DecoratorName, Value: "@resource", Line: 8, Column: 1},
		{Type: frontendTypes.DecoratorArg, Value: "aws.ecr.LifecyclePolicy", Line: 8, Column: 11},
		{Type: frontendTypes.DecoratorArg, Value: "ServerRepoLifecyclePolicy", Line: 8, Column: 36},
		{Type: frontendTypes.Newline, Value: "", Line: 8, Column: 62},
		{Type: frontendTypes.NaturalLanguage, Value: "In $ServerRepository", Line: 9, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 9, Column: 21},
		{
			Type:   frontendTypes.NaturalLanguage,
			Value:  "Policy retains only the last 10 untagged images in the repository. Images beyond this count will expire.",
			Line:   10,
			Column: 1,
		},
		{Type: frontendTypes.Newline, Value: "", Line: 10, Column: 105},
		{Type: frontendTypes.Newline, Value: "", Line: 11, Column: 1},
		{Type: frontendTypes.DecoratorName, Value: "@resource", Line: 12, Column: 1},
		{Type: frontendTypes.DecoratorArg, Value: "aws.ecs.Service", Line: 12, Column: 11},
		{Type: frontendTypes.DecoratorArg, Value: "ServerEcsService", Line: 12, Column: 28},
		{Type: frontendTypes.Newline, Value: "", Line: 12, Column: 45},
		{
			Type:   frontendTypes.NaturalLanguage,
			Value:  "In $EcsCluster, using $ServerTaskDefinition",
			Line:   13,
			Column: 1,
		},
		{Type: frontendTypes.Newline, Value: "", Line: 13, Column: 44},
		{Type: frontendTypes.NaturalLanguage, Value: "Name: server", Line: 14, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 14, Column: 13},
		{Type: frontendTypes.NaturalLanguage, Value: "Desired count: 1", Line: 15, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 15, Column: 17},
		{Type: frontendTypes.NaturalLanguage, Value: "Launch type: FARGATE", Line: 16, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 16, Column: 21},
		{Type: frontendTypes.NaturalLanguage, Value: "---", Line: 17, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 17, Column: 4},
		{Type: frontendTypes.NaturalLanguage, Value: "Network configuration:", Line: 18, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 18, Column: 23},
		{Type: frontendTypes.NaturalLanguage, Value: "  Assigned public IP", Line: 19, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 19, Column: 21},
		{Type: frontendTypes.NaturalLanguage, Value: "  Subnets: $PublicSubnetA and $PublicSubnetB", Line: 20, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 20, Column: 45},
		{Type: frontendTypes.NaturalLanguage, Value: "  Security group: $ServerEcsSecurityGroup", Line: 21, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 21, Column: 42},
		{Type: frontendTypes.NaturalLanguage, Value: "Load balancers:", Line: 22, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 22, Column: 16},
		{Type: frontendTypes.NaturalLanguage, Value: "  Target group: $ServerTargetGroup", Line: 23, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 23, Column: 35},
		{Type: frontendTypes.NaturalLanguage, Value: "  Container name: {server_container_name}", Line: 24, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 24, Column: 42},
		{Type: frontendTypes.NaturalLanguage, Value: "  Port: {container_port}", Line: 25, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 25, Column: 25},
		{Type: frontendTypes.NaturalLanguage, Value: "Deployment:", Line: 26, Column: 1},
		{Type: frontendTypes.Newline, Value: "", Line: 26, Column: 12},
		{
			Type:   frontendTypes.NaturalLanguage,
			Value:  "  ECS type deployment controller",
			Line:   27,
			Column: 1,
		},
		{Type: frontendTypes.Newline, Value: "", Line: 27, Column: 33},
		{
			Type:   frontendTypes.NaturalLanguage,
			Value:  "  Deployment circuit breaker: enabled with rollback",
			Line:   28,
			Column: 1,
		},
		{Type: frontendTypes.Newline, Value: "", Line: 28, Column: 52},
		{
			Type:   frontendTypes.NaturalLanguage,
			Value:  "  Wait for steady state: True",
			Line:   29,
			Column: 1,
		},
		{Type: frontendTypes.Newline, Value: "", Line: 29, Column: 30},
		{Type: frontendTypes.Newline, Value: "", Line: 30, Column: 1},
		{Type: frontendTypes.DecoratorName, Value: "@variable", Line: 31, Column: 1},
		{Type: frontendTypes.DecoratorArg, Value: "server_container_name", Line: 31, Column: 11},
		{Type: frontendTypes.DecoratorArg, Value: "string", Line: 31, Column: 34},
		{Type: frontendTypes.DecoratorArg, Value: "server-container", Line: 31, Column: 42},
		{Type: frontendTypes.Newline, Value: "", Line: 31, Column: 59},
		{Type: frontendTypes.NaturalLanguage, Value: "Description: Server container name", Line: 32, Column: 1},
		{Type: frontendTypes.EOF, Value: "", Line: 32, Column: 35},
	}
}
