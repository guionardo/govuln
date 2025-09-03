package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/guionardo/govuln/internal/config"
	"github.com/urfave/cli/v3"
)

var (
	storePath  string
	outputType string
)

func GetRoot() *cli.Command {
	cmd := &cli.Command{
		Version:     "v0.2.0",
		Description: "A comprehensive vulnerability scanner for Go projects with intelligent caching and submodule support",
		Usage:       "helper for golang.org/x/vuln/cmd/govulncheck",
		Commands: []*cli.Command{
			runCommand(),
			storeCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "store",
				Usage:       "Path of the store for caching checks of internal packages",
				Value:       config.Get().StoreDefaultPath,
				Destination: &storePath,
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "format of the output: color, markdown, json, yaml", //TODO: Implement this output formats
				Value: "color",
				Validator: func(v string) error {
					if !slices.Contains([]string{"color", "table", "markdown"}, v) {
						return fmt.Errorf("invalid output: %s", v)
					}
					return nil
				},
				Destination: &outputType,
			},
		},
	}

	return cmd
}

func Run() {
	cmd := GetRoot()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
