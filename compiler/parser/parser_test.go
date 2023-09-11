package parser_test

import (
	"salami/compiler/parser"
	"salami/compiler/types"
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

func getExpectedResources() []*types.Resource {
	return []*types.Resource{
		{
			ResourceType:        "cloudwatch.LogGroup",
			LogicalName:         "CumuliServerLogGroup",
			NaturalLanguage:     "Name: cumuli-server-log-group",
			Uses:                []types.LogicalName{},
			Exports:             make(map[string]string),
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
		{
			ResourceType:        "ecr.Repository",
			LogicalName:         "CumuliServerRepository",
			NaturalLanguage:     "Name: cumuli-server\nHas mutable image tags.",
			Uses:                []types.LogicalName{},
			Exports:             map[string]string{"name": "cumuli-server-ecr-repository-name"},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
		{
			ResourceType: "ecr.LifecyclePolicy",
			LogicalName:  "CumuliServerRepoLifecyclePolicy",
			NaturalLanguage: "Policy: A JSON policy with a rule that retains only the last 10 untagged images in the repository. " +
				"Images beyond this count will expire.",
			Uses:                []types.LogicalName{"CumuliServerRepository"},
			Exports:             make(map[string]string),
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
		{
			ResourceType: "aws.ecs.Service",
			LogicalName:  "CumuliServerEcsService",
			NaturalLanguage: `Name: cumuli-server
Desired count: 1
Launch type: FARGATE
---
Network configuration:
  Assigned public IP.
  Subnets: PublicSubnetA and PublicSubnetB.
  Security group: CumuliServerEcsSecurityGroup.
Load balancers:
  Target group: CumuliServerTargetGroup.
  Container name: {server_container_name}.
  Port: {container_port}.
Deployment:
  ECS type deployment controller.
  Deployment circuit breaker: enabled with rollback.
  Wait for steady state: True`,
			Uses: []types.LogicalName{
				"CumuliEcsCluster",
				"CumuliServerTaskDefinition",
				"PublicSubnetA",
				"PublicSubnetB",
				"CumuliServerEcsSecurityGroup",
				"CumuliServerTargetGroup",
			},
			Exports:             map[string]string{"name": "exported-name"},
			ReferencedVariables: []string{"server_container_name", "container_port"},
			SourceFilePath:      "dummy/file/path",
		},
	}
}

func getExpectedVariables(t *testing.T) []*types.Variable {
	variableType, err := types.StringToVariableType("string")
	if err != nil {
		t.Fatalf(err.Error())
	}
	return []*types.Variable{
		{
			Description:    "Server container name",
			Name:           "server_container_name",
			Default:        "cumuli-server-container",
			Type:           variableType,
			SourceFilePath: "dummy/file/path",
		},
	}
}

func getInput() []*types.Token {
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
		{Type: types.FieldName, Value: "Default", Line: 40, Column: 1},
		{Type: types.FieldValue, Value: "cumuli-server-container", Line: 40, Column: 8},
		{Type: types.EOF, Value: "", Line: 40, Column: 31},
	}
}
