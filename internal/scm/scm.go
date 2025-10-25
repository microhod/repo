package scm

import (
	"io"
	"net/url"
	"os"
	stdpath "path"
	"path/filepath"
	"strings"

	"github.com/microhod/repo/internal/path"
)

// Client is source code management system client e.g. git
type Client interface {
	// ParseRepoFromRemote parses a Repo object based on a raw remote URL
	ParseRepoFromRemote(rawURL string) (*Repo, error)
	// Clone pulls down the repo from the remote and stores it on `path`. This will also set `repo.Local` to `path`
	Clone(repo *Repo, path string, options *CloneOptions) error
	// FindRepos finds any repos in the current path (this will search recursively in all subfolders)
	FindRepos(path string) ([]*Repo, error)
}

type CloneOptions struct {
	Progress io.Writer
}

type Repo struct {
	Remote *url.URL
	Local  string
	Server string
	Owner  string
	Name   string
}

func (r *Repo) OrgaisedLocalPath(root string) string {
	return path.Clean(stdpath.Join(root, r.Server, r.Owner, r.Name))
}

func (r *Repo) IsOrganised(root string) bool {
	return strings.EqualFold(r.Local, r.OrgaisedLocalPath(root))
}

func (r *Repo) Organise(root string) error {
	if r.IsOrganised(root) {
		return nil
	}

	path := r.OrgaisedLocalPath(root)
	// make parent directories for new path
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	if err := os.Rename(r.Local, path); err != nil {
		return err
	}
	r.Local = path
	return nil
}
