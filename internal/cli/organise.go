package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/microhod/repo/internal/repo"
	"github.com/microhod/repo/internal/path"
	"github.com/microhod/repo/internal/terminal"
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
	err = terminal.WithSpinner("searching for repos...", func() (err error) {
		repos, err = app.client.FindRepos(basePath)
		return err
	})
	if err != nil {
		return fmt.Errorf("finding repos: %w", err)
	}

	moves := []Move{}
	for _, repo := range repos {
		organised := repo.LocalPath(app.cfg.Local.Root)
		// ignore repos which should not move (case insensitive)
		if !strings.EqualFold(organised, repo.Local) {
			moves = append(moves, Move{
				Repo:           repo,
				OrganisedLocal: organised,
			})
		}
	}
	if len(moves) < 1 {
		return nil
	}

	fmt.Println(Table(moves))

	fmt.Printf("move the repos as listed above? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		return nil
	}

	for _, move := range moves {
		if err := move.Do(); err != nil {
			return err
		}
	}

	fmt.Println("done 🎉")

	return nil
}

type Move struct {
	Repo           *repo.Repo
	OrganisedLocal string
}

func (move Move) Do() error {
	// make parent directories for new path
	if err := os.MkdirAll(filepath.Dir(move.OrganisedLocal), os.ModePerm); err != nil {
		return err
	}
	return os.Rename(move.Repo.Local, move.OrganisedLocal)
}

func Table(moves []Move) string {
	builder := &strings.Builder{}
	builder.WriteByte('\n')

	table := tablewriter.NewWriter(builder)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("➤")
	table.SetBorder(false)

	for _, m := range moves {
		table.Append([]string{
			path.CollapseHomeDir(m.Repo.Local),
			path.CollapseHomeDir(m.OrganisedLocal),
		})
	}

	table.Render()
	return builder.String()
}
