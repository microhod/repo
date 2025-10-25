package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"

	"github.com/microhod/repo/internal/path"
	"github.com/microhod/repo/internal/repo"
)

func (app *App) organise(ctx *cli.Context) error {
	var err error

	basePath := ctx.Args().First()
	if basePath == "" {
		basePath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("geting working directory: %w", err)
		}
	}

	var repos []*repo.Repo
	err = withSpinner("searching for repos...", func() (err error) {
		repos, err = app.client.FindRepos(basePath)
		return err
	})
	if err != nil {
		return fmt.Errorf("finding repos: %w", err)
	}

	var moves []*repo.Repo
	for _, repo := range repos {
		// ignore repos which should not move (case insensitive)
		if !repo.IsOrganised(app.cfg.Local.Root) {
			moves = append(moves, repo)
		}
	}
	if len(moves) == 0 {
		return nil
	}

	if !confirmMoves(app.cfg.Local.Root, moves) {
		return nil
	}
	for _, r := range moves {
		if err := r.Organise(app.cfg.Local.Root); err != nil {
			return err
		}
	}
	fmt.Println("done 🎉")
	return nil
}

func confirmMoves(root string, repo []*repo.Repo) bool {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, r := range repo {
		fmt.Fprintf(w, "%s\t>\t%s\t\n",
			path.CollapseHomeDir(r.Local),
			path.CollapseHomeDir(r.OrgaisedLocalPath(root)),
		)
	}
	w.Flush()

	fmt.Printf("move the repos as listed above? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)
	return confirm == "y"
}
