// Package scm deals with source code management systems
package scm

import (
	"io"

	"github.com/microhod/repo/internal/domain"
)

// Client is source code management system client e.g. git
type Client interface {
	// ParseRepoFromRemote parses a Repo object based on a raw remote URL
	ParseRepoFromRemote(rawURL string) (*domain.Repo, error)
	// Clone pulls down the repo from the remote and stores it on `path`. This will also set `repo.Local` to `path`
	Clone(repo *domain.Repo, path string, options *CloneOptions) error
	// FindRepos finds any repos in the current path (this will search recursively in all subfolders)
	FindRepos(path string) ([]*domain.Repo, error)
}

type CloneOptions struct {
	Progress io.Writer
}
