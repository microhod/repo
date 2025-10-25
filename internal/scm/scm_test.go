package scm

import (
	"errors"
	"os"
	"testing"
	"path"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo_OrganisedLocalPath(t *testing.T) {
	repo := Repo{
		Server: "example.com",
		Owner:  "bob",
		Name:   "builder",
	}
	assert.Equal(t, "/root/example.com/bob/builder", repo.OrgaisedLocalPath("/root"))
}

func TestRepo_IsOrganised(t *testing.T) {
	testCases := []struct {
		name     string
		repo     Repo
		expected bool
	}{
		{
			name: "already organised",
			repo: Repo{
				Local:  "/root/example.com/bob/builder",
				Server: "example.com",
				Owner:  "bob",
				Name:   "BUILDER", // shoud be case insensitive
			},
			expected: true,
		},
		{
			name: "already organised but under different root",
			repo: Repo{
				Local:  "/other/example.com/bob/builder",
				Server: "example.com",
				Owner:  "bob",
				Name:   "builder",
			},
			expected: false,
		},
		{
			name: "not organised",
			repo: Repo{
				Local:  "/other/path",
				Server: "example.com",
				Owner:  "bob",
				Name:   "builder",
			},
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.repo.IsOrganised("/root"))
		})
	}
}

func TestRepo_Move(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	testDir := path.Join(wd, "/testdata/tmp")
	// cleanup old runs
	err = os.RemoveAll(testDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("deleting [%s]: %s", testDir, err.Error())
	}
	require.NoError(t, os.MkdirAll(path.Join(testDir, "/unorganised"), os.ModePerm))

	repo := &Repo{
		Local:  path.Join(testDir, "/unorganised"),
		Server: "example.com",
		Owner:  "bob",
		Name:   "builder",
	}
	expectedPath := path.Join(testDir, "/example.com/bob/builder")

	// 1st organise
	err = repo.Organise(testDir)
	require.NoError(t, err, "first organise should not fail")
	require.Equal(t, expectedPath, repo.Local, "first organise should edit repo local path")
	_, err = os.Stat(expectedPath)
	require.NoError(t, err, "first organise should move folder") 

	// 2nd organise (shoud be no-op)
	err = repo.Organise(testDir)
	require.NoError(t, err, "second organise should not fail")
	require.Equal(t, expectedPath, repo.Local, "second organise should not edit repo local path")
	_, err = os.Stat(expectedPath)
	require.NoError(t, err, "second organise should not move folder") 
}
