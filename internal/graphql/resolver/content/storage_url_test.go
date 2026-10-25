package content

import (
	"context"
	"testing"

	"github.com/dictyBase/fp-go-loom/loomcheck"
	"github.com/dictyBase/fp-go-loom/pipelinecheck"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/stretchr/testify/require"
)

func newPBContent(raw string) *pb.Content {
	return &pb.Content{
		Data: &pb.ContentData{
			Id: 1,
			Attributes: &pb.ContentAttributes{
				Name:      "home",
				Slug:      "dsc-home",
				Namespace: "dsc",
				Content:   raw,
			},
		},
	}
}

func TestContentPipelineCheck(t *testing.T) {
	t.Parallel()
	pipelinecheck.Require(t, pipelinecheck.Config{
		Roots: []string{"."},
		IsEntrypoint: func(name string) bool {
			return name == "Content"
		},
		RequirePointFreeSeed:      true,
		RequirePointFreeBranching: true,
		RequireTrivialDelegate:    true,
	})
}

func TestContentLoomCheck(t *testing.T) {
	t.Parallel()
	loomcheck.Require(t, loomcheck.Config{Roots: []string{"."}})
}

func TestContentStorageURLNormalization(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		resolver *ContentResolver
		input    string
		expected string
	}{
		{
			name: "rewrites dev storage origin to prod origin",
			resolver: &ContentResolver{
				StorageURL: "https://slimestorage.dictycr.org",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"root":{"children":[{"type":"image","src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249-7b17-4ede-9e99-47aed2cbadd3"}]}}`,
			expected: `{"root":{"children":[{"type":"image","src":"https://slimestorage.dictycr.org/editor/assets/2024-10-25/333dc249-7b17-4ede-9e99-47aed2cbadd3"}]}}`,
		},
		{
			name: "rewrites prod storage origin to dev origin",
			resolver: &ContentResolver{
				StorageURL: "https://storage.dictybase.dev",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"root":{"children":[{"type":"image","src":"https://slimestorage.dictycr.org/editor/assets/2026-09-22/0b9876d2-8d5b-47db-b21b-b08b9a0b8524"}]}}`,
			expected: `{"root":{"children":[{"type":"image","src":"https://storage.dictybase.dev/editor/assets/2026-09-22/0b9876d2-8d5b-47db-b21b-b08b9a0b8524"}]}}`,
		},
		{
			name: "identical origin with trailing slash is a no-op",
			resolver: &ContentResolver{
				StorageURL: "https://storage.dictybase.dev/",
				Bucket:     "/editor/",
				BucketPath: "/assets/",
			},
			input:    `{"root":{"children":[{"type":"image","src":"https://storage.dictybase.dev/editor/assets/2026-09-22/0b9876d2-8d5b-47db-b21b-b08b9a0b8524"}]}}`,
			expected: `{"root":{"children":[{"type":"image","src":"https://storage.dictybase.dev/editor/assets/2026-09-22/0b9876d2-8d5b-47db-b21b-b08b9a0b8524"}]}}`,
		},
		{
			name: "rewrites multiple storage URLs and preserves external URLs",
			resolver: &ContentResolver{
				StorageURL: "https://slimestorage.dictycr.org",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"a":"https://storage.dictybase.dev/editor/assets/2024-10-25/first","b":"https://dictybase.org/papers/123","c":"https://storage.dictybase.dev/other-bucket/file.png","d":"https://storage.dictybase.dev/editor/assets/2024-10-25/second"}`,
			expected: `{"a":"https://slimestorage.dictycr.org/editor/assets/2024-10-25/first","b":"https://dictybase.org/papers/123","c":"https://storage.dictybase.dev/other-bucket/file.png","d":"https://slimestorage.dictycr.org/editor/assets/2024-10-25/second"}`,
		},
		{
			name: "leaves external URLs containing the bucket path deeper in the path",
			resolver: &ContentResolver{
				StorageURL: "https://slimestorage.dictycr.org",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"a":"https://doi.org/10.1000/editor/assets/paper","b":"https://github.com/dictybase/editor/assets/code"}`,
			expected: `{"a":"https://doi.org/10.1000/editor/assets/paper","b":"https://github.com/dictybase/editor/assets/code"}`,
		},
		{
			name: "rewrites only the first bucket prefix of a repeated path",
			resolver: &ContentResolver{
				StorageURL: "https://slimestorage.dictycr.org",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"src":"https://storage.dictybase.dev/editor/assets/sub/editor/assets/uuid"}`,
			expected: `{"src":"https://slimestorage.dictycr.org/editor/assets/sub/editor/assets/uuid"}`,
		},
		{
			name: "returns input unchanged when storage config is incomplete",
			resolver: &ContentResolver{
				StorageURL: "",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249"}`,
			expected: `{"src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249"}`,
		},
		{
			name: "returns input unchanged when storage URL has only slashes",
			resolver: &ContentResolver{
				StorageURL: "///",
				Bucket:     "editor",
				BucketPath: "assets",
			},
			input:    `{"src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249"}`,
			expected: `{"src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249"}`,
		},
		{
			name: "returns input unchanged when bucket contains invalid UTF-8",
			resolver: &ContentResolver{
				StorageURL: "https://slimestorage.dictycr.org",
				Bucket:     string([]byte{0xff, 0xfe}),
				BucketPath: "assets",
			},
			input:    `{"src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249"}`,
			expected: `{"src":"https://storage.dictybase.dev/editor/assets/2024-10-25/333dc249"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tc.resolver.Content(
				context.Background(),
				newPBContent(tc.input),
			)
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
