package main

import (
	"fmt"
	"os"

	"github.com/microhod/repo/internal/cli"
	"github.com/microhod/repo/internal/config"
)

func main() {
	// parse config
	cfg, err := config.Parse()
	if err != nil {
		fmt.Printf("ERROR: failed to parse config: %s\n", err.Error())
		os.Exit(1)
	}
	// run cli app
	if err := cli.NewApp(cfg).Run(os.Args); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
