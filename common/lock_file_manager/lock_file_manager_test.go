package lock_file_manager_test

import (
	"path/filepath"
	"salami/common/lock_file_manager"
	"salami/common/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTargetFilesMeta(t *testing.T) {
	t.Run("should get target files meta from the lock file", func(t *testing.T) {
		setLockFile(t, "valid.toml")
		lock_file_manager.ValidateLockFile()
		expectedTargetFilesMeta := []types.TargetFileMeta{
			{FilePath: "path/to/target_file_1", Checksum: "e460d56360c0c4d1ff32fd5e5a56eb99"},
			{FilePath: "path/to/target_file_2", Checksum: "ea441ff260d926a935cf47abf698482d"},
		}
		actualTargetFilesMeta := lock_file_manager.GetTargetFilesMeta()
		require.ElementsMatch(
			t,
			expectedTargetFilesMeta,
			actualTargetFilesMeta,
		)
	})
}

func TestGetObjects(t *testing.T) {
	t.Run("should get objects from the lock file", func(t *testing.T) {
		setLockFile(t, "valid.toml")
		lock_file_manager.ValidateLockFile()
		expectedObjects := getExpectedObjects()
		actualObjects := lock_file_manager.GetObjects()
		for i, actualObject := range actualObjects {
			require.Equal(t, expectedObjects[i], actualObject)
		}
	})
}

func getExpectedObjects() []*types.Object {
	return []*types.Object{
		{
			SourceFilePath: "path/to/source_file",
			ParsedResource: &types.ParsedResource{
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
					Content: "resource \"aws_s3_bucket\" \"AssumedRolesBucket\" {\n" +
						"  bucket = \"assumed-roles\"\n  versioning {\n    enabled = true\n  }\n}",
				},
				{
					SegmentType:    types.CodeSegmentType("Export"),
					Content:        "output \"assumed-roles-bucket-name\" {\n  value = aws_s3_bucket.AssumedRolesBucket.bucket\n}",
				},
			},
		},
		{
			SourceFilePath: "path/to/source_file",
			ParsedResource: &types.ParsedResource{
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
					Content: "resource \"aws_s3_bucket_public_access_block\" \"AssetsPublicAccessBlock\" {\n" +
						"  bucket = aws_s3_bucket.AssumedRolesBucket.id\n\n  block_public_acls       = true\n" +
						"  block_public_policy     = false\n  ignore_public_acls      = true\n" +
						"  restrict_public_buckets = false\n}",
				},
			},
		},
		{
			SourceFilePath: "path/to/source_file",
			ParsedResource: &types.ParsedResource{
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
			ParsedVariable: &types.ParsedVariable{
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
					Content: "variable \"server_container_name\" {\n  description = \"Server container name\"\n" +
						"  type        = string\n  default     = \"server-container\"\n}",
				},
			},
		},
	}
}

func setLockFile(t *testing.T, fileName string) {
	fixturePath := filepath.Join("testdata", "lock_files", fileName)
	lock_file_manager.SetLockFilePath(fixturePath)
}
