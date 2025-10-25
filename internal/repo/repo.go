package repo

import (
	"net/url"
	"path/filepath"
	stdpath "path"
	"os"
	"strings"

	"github.com/microhod/repo/internal/path"
)

type Repo struct {
	Remote *url.URL
	Local  string
	Server string
	Owner  string
	Name   string
}

func (r Repo) OrgaisedLocalPath(root string) string {
	return path.Clean(stdpath.Join(root, r.Server, r.Owner, r.Name))
}

func (r Repo) IsOrganised(root string) bool {
	return strings.EqualFold(r.Local, r.OrgaisedLocalPath(root))
}


func (r Repo) Organise(root string) error {
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
