package commands

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/guionardo/govuln/internal/params"
	"github.com/urfave/cli/v3"
)

var (
	storePath  string
	outputType string
)

func outputValidator(v string) error {
	if !slices.Contains([]string{"color", "table", "markdown"}, v) {
		return fmt.Errorf("invalid output: %s", v)
	}
	return nil
}
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
				Value:       params.STORE_DEFAULT_PATH,
				Destination: &storePath,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       "format of the output: color, markdown, json, yaml", //TODO: Implement this output formats
				Value:       "color",
				Validator:   outputValidator,
				Destination: &outputType,
			},
		},
	}

	return cmd
}

func Run() int {
	cmd := GetRoot()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		return 1
	}
	return 0
}
