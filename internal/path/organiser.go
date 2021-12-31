package path

import (
	"path"

	"github.com/microhod/repo/internal/domain"
)

type Organiser interface {
	// Organise finds the path to in which to organise the current repo
	Organise(repo *domain.Repo) string
}

type SimpleOrganiser struct {
	Root string
}

func (o *SimpleOrganiser) Organise(repo *domain.Repo) string {
	return Clean(path.Join(o.Root, repo.Server, repo.Owner, repo.Name))
}
