package change_manager_test

import (
	"salami/common/change_manager"
	"salami/common/symbol_table"
	"salami/common/types"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateChangeSet(t *testing.T) {
	t.Run("should return an empty change set when there are no changes", func(t *testing.T) {
		previousObjects, previousResources, previousVariables := getPreviousObjects()
		symbolTable, err := symbol_table.NewSymbolTable(previousResources, previousVariables)
		require.NoError(t, err)
		changeSet := change_manager.GenerateChangeSet(previousObjects, symbolTable)
		require.Equal(t, changeSet, &types.ChangeSet{Diffs: []types.ChangeSetDiff{}})
	})

	t.Run("should return a change set with additions, deletions, and changes when they exist", func(t *testing.T) {
		previousObjects, previousResources, previousVariables := getPreviousObjects()
		newResources := previousResources[:1]
		newResources = append(
			newResources,
			&types.ParsedResource{
				ResourceType:        types.ResourceType("aws.s3.BucketPolicy"),
				LogicalName:         types.LogicalName("AssumedRolesBucketPolicy"),
				NaturalLanguage:     "Policy: A JSON policy that denies everyone access",
				Uses:                []types.LogicalName{"AssumedRolesBucket"},
				Exports:             map[string]string{},
				ReferencedVariables: []string{},
				SourceFilePath:      "path/to/source_file",
				SourceFileLine:      0,
			},
			&types.ParsedResource{
				ResourceType:        types.ResourceType("aws.logs.LogGroup"),
				LogicalName:         types.LogicalName("LogsGroup"),
				NaturalLanguage:     "Log group: /ecs/my-app\nRetention: 30 days",
				Uses:                []types.LogicalName{},
				Exports:             map[string]string{},
				ReferencedVariables: []string{},
				SourceFilePath:      "path/to/source_file",
				SourceFileLine:      1,
			},
		)
		previousVariables = append(previousVariables, &types.ParsedVariable{
			Name:           "container port",
			Description:    "Container port",
			Type:           types.VariableType("string"),
			Default:        "8080",
			SourceFilePath: "path/to/source_file",
			SourceFileLine: 15,
		})
		symbolTable, err := symbol_table.NewSymbolTable(newResources, previousVariables)
		require.NoError(t, err)
		changeSet := change_manager.GenerateChangeSet(previousObjects, symbolTable)
		changeSetDiffs := sortChangeSetDiffs(changeSet.Diffs)
		require.Equal(t, 4, len(changeSetDiffs))
		expectedDiffs := getExpectedChangeSetDiffs()
		for i, actualDiff := range changeSetDiffs {
			require.Equal(t, expectedDiffs[i], actualDiff)
		}
	})
}

func getPreviousObjects() ([]*types.Object, []*types.ParsedResource, []*types.ParsedVariable) {
	objects := []*types.Object{
		{
			SourceFilePath: "path/to/source_file",
			Parsed: &types.ParsedResource{
				ResourceType:        types.ResourceType("aws.s3.Bucket"),
				LogicalName:         types.LogicalName("AssumedRolesBucket"),
				NaturalLanguage:     "Bucket: assumed-roles\nVersioning enabled: True",
				Uses:                []types.LogicalName{},
				Exports:             map[string]string{"name": "assumed-roles-bucket-name"},
				ReferencedVariables: []string{},
				SourceFilePath:      "path/to/source_file",
				SourceFileLine:      0,
			},
			CodeSegments: []types.CodeSegment{
				{
					SegmentType:    types.CodeSegmentType("Resource"),
					TargetFilePath: "path/to/target_file_1",
					Content: "resource \"aws_s3_bucket\" \"AssumedRolesBucket\" {\n" +
						"  bucket = \"assumed-roles\"\n  versioning {\n    enabled = true\n  }\n}",
				},
				{
					SegmentType:    types.CodeSegmentType("Export"),
					TargetFilePath: "path/to/target_file_1",
					Content:        "output \"assumed-roles-bucket-name\" {\n  value = aws_s3_bucket.AssumedRolesBucket.bucket\n}",
				},
			},
		},
		{
			SourceFilePath: "path/to/source_file",
			Parsed: &types.ParsedResource{
				ResourceType: types.ResourceType("aws.s3.BucketPublicAccessBlock"),
				LogicalName:  types.LogicalName("AssetsPublicAccessBlock"),
				NaturalLanguage: "Block public ACLs: True\nBlock public policy: False\n" +
					"Ignore public ACLs: True\nRestrict public buckets: False",
				Uses:                []types.LogicalName{"AssumedRolesBucket"},
				Exports:             map[string]string{},
				ReferencedVariables: []string{},
				SourceFilePath:      "path/to/source_file",
				SourceFileLine:      0,
			},
			CodeSegments: []types.CodeSegment{
				{
					SegmentType:    types.CodeSegmentType("Resource"),
					TargetFilePath: "path/to/target_file_1",
					Content: "resource \"aws_s3_bucket_public_access_block\" \"AssetsPublicAccessBlock\" {\n" +
						"  bucket = aws_s3_bucket.AssumedRolesBucket.id\n\n  block_public_acls       = true\n" +
						"  block_public_policy     = false\n  ignore_public_acls      = true\n" +
						"  restrict_public_buckets = false\n}",
				},
			},
		},
		{
			SourceFilePath: "path/to/source_file",
			Parsed: &types.ParsedResource{
				ResourceType: types.ResourceType("aws.s3.BucketPolicy"),
				LogicalName:  types.LogicalName("AssumedRolesBucketPolicy"),
				NaturalLanguage: "Policy: A JSON policy that allows all principals to perform the " +
					"\"s3:GetObject\" action on all objects in the specified S3 bucket.",
				Uses:                []types.LogicalName{"AssumedRolesBucket"},
				Exports:             map[string]string{},
				ReferencedVariables: []string{},
				SourceFilePath:      "path/to/source_file",
				SourceFileLine:      0,
			},
			CodeSegments: []types.CodeSegment{
				{
					SegmentType:    types.CodeSegmentType("Resource"),
					TargetFilePath: "path/to/target_file_1",
					Content: "resource \"aws_s3_bucket_policy\" \"AssumedRolesBucketPolicy\" {\n" +
						"  bucket = aws_s3_bucket.AssumedRolesBucket.id\n\n  policy = jsonencode({\n" +
						"    Version = \"2012-10-17\"\n    Statement = [\n      {\n" +
						"        Action   = \"s3:GetObject\"\n        Effect   = \"Allow\"\n" +
						"        Resource = \"${aws_s3_bucket.AssumedRolesBucket.arn}/*\"\n" +
						"        Principal = \"*\"\n      }\n    ]\n  })\n}",
				},
			},
		},
		{
			SourceFilePath: "path/to/source_file",
			Parsed: &types.ParsedVariable{
				Name:           "server_container_name",
				Description:    "Server container name",
				Type:           types.VariableType("string"),
				Default:        "server-container",
				SourceFilePath: "path/to/source_file",
				SourceFileLine: 0,
			},
			CodeSegments: []types.CodeSegment{
				{
					SegmentType:    types.CodeSegmentType("Variable"),
					TargetFilePath: "path/to/target_file_2",
					Content: "variable \"server_container_name\" {\n  description = \"Server container name\"\n" +
						"  type        = string\n  default     = \"server-container\"\n}",
				},
			},
		},
	}

	resources := []*types.ParsedResource{}
	variables := []*types.ParsedVariable{}
	for _, obj := range objects {
		if res, ok := obj.Parsed.(*types.ParsedResource); ok {
			resources = append(resources, res)
		} else if var_, ok := obj.Parsed.(*types.ParsedVariable); ok {
			variables = append(variables, var_)
		}
	}
	return objects, resources, variables
}

func getExpectedChangeSetDiffs() []types.ChangeSetDiff {
	return []types.ChangeSetDiff{
		{
			OldObject: nil,
			NewObject: &types.Object{
				SourceFilePath: "path/to/source_file",
				Parsed: &types.ParsedVariable{
					Name:           "container port",
					Description:    "Container port",
					Type:           types.VariableType("string"),
					Default:        "8080",
					SourceFilePath: "path/to/source_file",
					SourceFileLine: 15,
				},
			},
		},
		{
			OldObject: nil,
			NewObject: &types.Object{
				SourceFilePath: "path/to/source_file",
				Parsed: &types.ParsedResource{
					ResourceType:        types.ResourceType("aws.logs.LogGroup"),
					LogicalName:         types.LogicalName("LogsGroup"),
					NaturalLanguage:     "Log group: /ecs/my-app\nRetention: 30 days",
					Uses:                []types.LogicalName{},
					Exports:             map[string]string{},
					ReferencedVariables: []string{},
					SourceFilePath:      "path/to/source_file",
					SourceFileLine:      1,
				},
			},
		},
		{
			OldObject: &types.Object{
				SourceFilePath: "path/to/source_file",
				Parsed: &types.ParsedResource{
					ResourceType: types.ResourceType("aws.s3.BucketPublicAccessBlock"),
					LogicalName:  types.LogicalName("AssetsPublicAccessBlock"),
					NaturalLanguage: "Block public ACLs: True\nBlock public policy: False\n" +
						"Ignore public ACLs: True\nRestrict public buckets: False",
					Uses:                []types.LogicalName{"AssumedRolesBucket"},
					Exports:             map[string]string{},
					ReferencedVariables: []string{},
					SourceFilePath:      "path/to/source_file",
					SourceFileLine:      0,
				},
				CodeSegments: []types.CodeSegment{
					{
						SegmentType:    types.CodeSegmentType("Resource"),
						TargetFilePath: "path/to/target_file_1",
						Content: "resource \"aws_s3_bucket_public_access_block\" \"AssetsPublicAccessBlock\" {\n" +
							"  bucket = aws_s3_bucket.AssumedRolesBucket.id\n\n  block_public_acls       = true\n" +
							"  block_public_policy     = false\n  ignore_public_acls      = true\n" +
							"  restrict_public_buckets = false\n}",
					},
				},
			},
			NewObject: nil,
		},
		{
			OldObject: &types.Object{
				SourceFilePath: "path/to/source_file",
				Parsed: &types.ParsedResource{
					ResourceType: types.ResourceType("aws.s3.BucketPolicy"),
					LogicalName:  types.LogicalName("AssumedRolesBucketPolicy"),
					NaturalLanguage: "Policy: A JSON policy that allows all principals to perform the " +
						"\"s3:GetObject\" action on all objects in the specified S3 bucket.",
					Uses:                []types.LogicalName{"AssumedRolesBucket"},
					Exports:             map[string]string{},
					ReferencedVariables: []string{},
					SourceFilePath:      "path/to/source_file",
					SourceFileLine:      0,
				},
				CodeSegments: []types.CodeSegment{
					{
						SegmentType:    types.CodeSegmentType("Resource"),
						TargetFilePath: "path/to/target_file_1",
						Content: "resource \"aws_s3_bucket_policy\" \"AssumedRolesBucketPolicy\" {\n" +
							"  bucket = aws_s3_bucket.AssumedRolesBucket.id\n\n  policy = jsonencode({\n" +
							"    Version = \"2012-10-17\"\n    Statement = [\n      {\n" +
							"        Action   = \"s3:GetObject\"\n        Effect   = \"Allow\"\n" +
							"        Resource = \"${aws_s3_bucket.AssumedRolesBucket.arn}/*\"\n" +
							"        Principal = \"*\"\n      }\n    ]\n  })\n}",
					},
				},
			},
			NewObject: &types.Object{
				SourceFilePath: "path/to/source_file",
				Parsed: &types.ParsedResource{
					ResourceType:        types.ResourceType("aws.s3.BucketPolicy"),
					LogicalName:         types.LogicalName("AssumedRolesBucketPolicy"),
					NaturalLanguage:     "Policy: A JSON policy that denies everyone access",
					Uses:                []types.LogicalName{"AssumedRolesBucket"},
					Exports:             map[string]string{},
					ReferencedVariables: []string{},
					SourceFilePath:      "path/to/source_file",
					SourceFileLine:      0,
				},
			},
		},
	}
}

func sortChangeSetDiffs(diffs []types.ChangeSetDiff) []types.ChangeSetDiff {
	sort.Slice(diffs, func(i, j int) bool {
		iDiff := diffs[i]
		jDiff := diffs[j]

		iOldIsVar := false
		iNewIsVar := false
		jOldIsVar := false
		jNewIsVar := false

		iOldName := ""
		iNewName := ""
		jOldName := ""
		jNewName := ""

		if iDiff.OldObject != nil {
			if _, iOldIsVar = iDiff.OldObject.Parsed.(*types.ParsedVariable); iOldIsVar {
				iOldName = iDiff.OldObject.Parsed.(*types.ParsedVariable).Name
			} else {
				iOldName = string(iDiff.OldObject.Parsed.(*types.ParsedResource).LogicalName)
			}
		}

		if iDiff.NewObject != nil {
			if _, iNewIsVar = iDiff.NewObject.Parsed.(*types.ParsedVariable); iNewIsVar {
				iNewName = iDiff.NewObject.Parsed.(*types.ParsedVariable).Name
			} else {
				iNewName = string(iDiff.NewObject.Parsed.(*types.ParsedResource).LogicalName)
			}
		}

		if jDiff.OldObject != nil {
			if _, jOldIsVar = jDiff.OldObject.Parsed.(*types.ParsedVariable); jOldIsVar {
				jOldName = jDiff.OldObject.Parsed.(*types.ParsedVariable).Name
			} else {
				jOldName = string(jDiff.OldObject.Parsed.(*types.ParsedResource).LogicalName)
			}
		}

		if jDiff.NewObject != nil {
			if _, jNewIsVar = jDiff.NewObject.Parsed.(*types.ParsedVariable); jNewIsVar {
				jNewName = jDiff.NewObject.Parsed.(*types.ParsedVariable).Name
			} else {
				jNewName = string(jDiff.NewObject.Parsed.(*types.ParsedResource).LogicalName)
			}
		}

		if iOldIsVar != jOldIsVar {
			return iOldIsVar
		}

		if iNewIsVar != jNewIsVar {
			return iNewIsVar
		}

		if iOldName != jOldName {
			return iOldName < jOldName
		}

		return iNewName < jNewName
	})
	return diffs
}
