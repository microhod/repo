package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
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
		return fmt.Errorf("parsing repo: %w", err)
	}
	path := repo.OrgaisedLocalPath(app.cfg.Local.Root)

	// clone
	err = withSpinner("cloning...", func() error {
		return app.client.Clone(repo, path, nil)
	})
	if err != nil {
		return fmt.Errorf("cloning repo: %w", err)
	}

	fmt.Println(repo.Local)
	return nil
}

func withSpinner(message string, f func() error) error {
	// create basic spinner coloured spinner
	s := spinner.New(
		[]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
	)
	s.Color("blue")
	s.Suffix = " " + message
	// clear terminal line after spinner is stopped
	s.FinalMSG = "\033[2K\r"

	s.Start()
	defer s.Stop()
	return f()
}
