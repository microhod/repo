// Package cli describes the command line application
package cli

import (
	"github.com/urfave/cli/v2"

	"github.com/microhod/repo/internal/config"
	"github.com/microhod/repo/internal/path"
	"github.com/microhod/repo/internal/scm"
)

type App struct {
	client    scm.Client
	organiser path.Organiser
}

func NewApp(cfg config.Config) *cli.App {
	app := &App{
		organiser: &path.SimpleOrganiser{
			Root: cfg.Local.Root,
		},
		client: &scm.Git{
			DefaultRemotePrefix: cfg.Remote.Default.Prefix,
		},
	}
	cliApp := &cli.App{
		Usage: "A cli application to organise scm repositories in a structured hierarchy",
		Commands: []*cli.Command{
			{
				Name:   "clone",
				Usage:  "clone a repo",
				Action: app.clone,
			},
			{
				Name:   "organise",
				Usage:  "organise all repos under the current path into a structured heirachy",
				Action: app.organise,
			},
		},
	}
	return cliApp
}
