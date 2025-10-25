package path

import (
	"os"
	"path"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "standard cleaning",
			path:     "/base/path/to/../../something",
			expected: "/base/something",
		},
		{
			name:     "with ~ prefix",
			path:     "~/base/path/../something",
			expected: inHomeDir("/base/something"),
		},
		{
			name:     "with ~ non-prefix",
			path:     "/base/path/../~/something",
			expected: "/base/~/something",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, Clean(tc.path))
		})
	}
}

func TestCollapseHomeDir(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "has home dir prefix",
			path:     inHomeDir("/some/path"),
			expected: "~/some/path",
		},
		{
			name:     "no home dir prefix",
			path:     "/some/path",
			expected: "/some/path",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, CollapseHomeDir(tc.path))
		})
	}
}

func TestFindDir(t *testing.T) {
	/*
	testdir
	├── a
	│   ├── aa
	│   │   └── findme  <- nested matches should be ignored
	│   └── findme
	├── b
	│   └── findme
	└── c.txt
	*/
	wd, err := os.Getwd()
	require.NoError(t, err)
	expected := []string{
		path.Join(wd, "testdir/a/findme"),
		path.Join(wd, "testdir/b/findme"),
	}

	paths, err := FindDir(path.Join(wd, "testdir"), "findme")
	require.NoError(t, err)

	sort.Strings(expected)
	sort.Strings(paths)
	assert.Equal(t, expected, paths)
}

func inHomeDir(p string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(home, p)
}
