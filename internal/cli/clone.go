package cli

import (
	"fmt"

	"github.com/microhod/repo/internal/terminal"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *App) clone(ctx *cli.Context) error {
	if ctx.Args().Len() < 1 {
		return fmt.Errorf("no repo remote specified")
	}
	rawRemoteURL := ctx.Args().First()

	// parse repo
	repo, err := app.client.ParseRepoFromRemote(rawRemoteURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse repo")
	}
	path := app.organiser.Organise(repo)

	// clone
	err = terminal.WithSpinner("cloning...", func() error {
		return app.client.Clone(repo, path, nil)
	})
	if err != nil {
		return errors.Wrap(err, "failed to clone repo")
	}

	fmt.Println(repo.Local)

	return nil
}
