package repo

import (
	"net/url"
	stdpath "path"

	"github.com/microhod/repo/internal/path"
)

type Repo struct {
	Remote *url.URL
	Local  string
	Server string
	Owner  string
	Name   string
}

func (r Repo) LocalPath(root string) string {
	return path.Clean(stdpath.Join(root, r.Server, r.Owner, r.Name))
}
