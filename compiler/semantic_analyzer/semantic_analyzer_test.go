package semantic_analyzer_test

import (
	"salami/compiler/semantic_analyzer"
	"salami/compiler/symbol_table"
	"salami/compiler/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceRequiredFields(t *testing.T) {
	resources := []*types.Resource{
		{
			ResourceType:        "",
			LogicalName:         "CumuliServerLogGroup",
			NaturalLanguage:     "Name: cumuli-server-log-group",
			Uses:                []types.LogicalName{},
			Exports:             make(map[string]string),
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
	}
	semanticAnalyzer := createSemanticAnalyzer(t, resources, []*types.Variable{})
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

	resources = []*types.Resource{
		{
			ResourceType:        "cloudwatch.LogGroup",
			NaturalLanguage:     "Name: cumuli-server-log-group",
			Uses:                []types.LogicalName{},
			Exports:             make(map[string]string),
			ReferencedVariables: []string{},
			SourceFilePath:      "dummy/file/path",
		},
	}
	semanticAnalyzer = createSemanticAnalyzer(t, resources, []*types.Variable{})
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
	variables := []*types.Variable{
		{
			Name:           "",
			Description:    "Test variable",
			Value:          "test-value",
			SourceFilePath: "dummy/file/path",
		},
	}
	semanticAnalyzer := createSemanticAnalyzer(t, []*types.Resource{}, variables)
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
	resources := []*types.Resource{
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
	variables := []*types.Variable{
		{
			Name:           "server_container_name",
			Description:    "Name of the container that runs the server",
			Value:          "cumuli-server",
			SourceFilePath: "dummy/file/path",
		},
	}

	semanticAnalyzer := createSemanticAnalyzer(t, resources, variables)
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: Referenced variable 'container_port' is not defined"
		require.Equal(
			t,
			expectedErrorMessage,
			err.Error(),
			"Expected error message '%s' but got '%s'", expectedErrorMessage, err.Error(),
		)
	}
}

func TestUsedResourcesExist(t *testing.T) {
	resources := []*types.Resource{
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
	}
	semanticAnalyzer := createSemanticAnalyzer(t, resources, []*types.Variable{})
	if err := semanticAnalyzer.Analyze(); err != nil {
		require.NotNil(t, err, "Expected error but got nil")
		expectedErrorMessage := "\ndummy/file/path\n  semantic error: Used resource 'CumuliServerRepository' is not defined"
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
	resources []*types.Resource,
	variables []*types.Variable,
) *semantic_analyzer.SemanticAnalyzer {
	symbolTable, err := symbol_table.NewSymbolTable(resources, variables)
	if err != nil {
		t.Fatal(err)
	}
	return semantic_analyzer.NewSemanticAnalyzer(symbolTable)
}
