package semantic_analyzer_test

import (
	"salami/common/symbol_table"
	"salami/common/types"
	"salami/frontend/semantic_analyzer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceRequiredFields(t *testing.T) {
	resources := []*types.ParsedResource{
		{
			ResourceType:        "",
			LogicalName:         "ServerLogGroup",
			NaturalLanguage:     "Name: server-log-group",
			ReferencedResources: []types.LogicalName{},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
	}
	semanticAnalyzer := createSemanticAnalyzer(t, resources, []*types.ParsedVariable{})
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: " +
			"Resource type field on a resource object is missing or empty"
		require.Equal(
			t,
			expectedErrorMessage,
			err.Error(),
			"Expected error message '%s' but got '%s'",
			expectedErrorMessage,
			err.Error(),
		)
	}

	resources = []*types.ParsedResource{
		{
			ResourceType:        "cloudwatch.LogGroup",
			NaturalLanguage:     "Name: server-log-group",
			ReferencedResources: []types.LogicalName{},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
	}
	semanticAnalyzer = createSemanticAnalyzer(t, resources, []*types.ParsedVariable{})
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: " +
			"Logical name field on a resource object is missing or empty"
		require.Equal(
			t,
			expectedErrorMessage,
			err.Error(),
			"Expected error message '%s' but got '%s'",
			expectedErrorMessage,
			err.Error(),
		)
	}
}

func TestVariableRequiredFields(t *testing.T) {
	variables := []*types.ParsedVariable{
		{
			Name:            "",
			NaturalLanguage: "Test variable",
			Default:         "test-value",
			SourceFilePath:  "dummy/file/path",
		},
	}
	semanticAnalyzer := createSemanticAnalyzer(t, []*types.ParsedResource{}, variables)
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: " +
			"Name field on a variable object is missing or empty"
		require.Equal(
			t,
			expectedErrorMessage,
			err.Error(),
			"Expected error message '%s' but got '%s'",
			expectedErrorMessage,
			err.Error(),
		)
	}
}

func TestReferencedVariablesAreDefined(t *testing.T) {
	resources := []*types.ParsedResource{
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
			ReferencedResources: []types.LogicalName{
				"EcsCluster",
				"ServerTaskDefinition",
				"PublicSubnetA",
				"PublicSubnetB",
				"ServerEcsSecurityGroup",
				"ServerTargetGroup",
			},
			ReferencedVariables: []string{"server_container_name", "container_port"},
			SourceFilePath:      "dummy/file/path",
		},
	}
	variables := []*types.ParsedVariable{
		{
			Name:            "server_container_name",
			NaturalLanguage: "Name of the container that runs the server",
			Default:         "server",
			SourceFilePath:  "dummy/file/path",
		},
	}

	semanticAnalyzer := createSemanticAnalyzer(t, resources, variables)
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: referenced variable 'container_port' is not defined"
		require.Equal(
			t,
			expectedErrorMessage,
			err.Error(),
			"Expected error message '%s' but got '%s'", expectedErrorMessage, err.Error(),
		)
	}
}

func TestUsedResourcesExist(t *testing.T) {
	resources := []*types.ParsedResource{
		{
			ResourceType:        "ecr.LifecyclePolicy",
			LogicalName:         "ServerRepoLifecyclePolicy",
			NaturalLanguage:     "Policy retains only the last 10 untagged images in the repository. Images beyond this count will expire.",
			ReferencedResources: []types.LogicalName{"ServerRepository"},
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
	}
	semanticAnalyzer := createSemanticAnalyzer(t, resources, []*types.ParsedVariable{})
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: referenced resource 'ServerRepository' is not defined"
		require.Equal(
			t,
			expectedErrorMessage,
			err.Error(),
			"Expected error message '%s' but got '%s'", expectedErrorMessage, err.Error(),
		)
	}
}

func createSemanticAnalyzer(
	t *testing.T,
	resources []*types.ParsedResource,
	variables []*types.ParsedVariable,
) *semantic_analyzer.SemanticAnalyzer {
	symbolTable, err := symbol_table.NewSymbolTable(resources, variables)
	if err != nil {
		t.Fatal(err)
	}
	return semantic_analyzer.NewSemanticAnalyzer(symbolTable)
}
