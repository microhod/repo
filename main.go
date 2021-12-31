package main

import (
	"log"
	"os"

	"github.com/microhod/repo/internal/cli"
	"github.com/microhod/repo/internal/config"
)

func main() {
	// disable all log formatting e.g. timestamps
	log.SetFlags(0)
	// parse config
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("failed to parse config: %s\n", err.Error())
	}
	// run cli app
	if err := cli.NewApp(cfg).Run(os.Args); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
