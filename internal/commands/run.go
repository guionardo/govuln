package commands

import (
	"context"
	"fmt"

	"github.com/guionardo/govuln/internal/check"
	"github.com/guionardo/govuln/internal/params"
	"github.com/guionardo/govuln/internal/store"

	pathtools "github.com/guionardo/go/pkg/path_tools"
	"github.com/urfave/cli/v3"
)

var (
	projectPath   string
	justWarn      bool
	dontCheckSubs bool
	internalOwner string
)

func runCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "analyses the project throught govulncheck",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Usage:       "Path of the project to check vulnerabilities",
				Value:       params.CURRENT_PATH,
				Destination: &projectPath,
			},
			&cli.BoolFlag{
				Name:        "just-warn",
				Value:       false,
				Usage:       "Just warn about vulnerabilities without failing",
				Destination: &justWarn,
			},
			&cli.BoolFlag{
				Name:        "dont-check-subs",
				Value:       false,
				Usage:       "Don't check submodules",
				Destination: &dontCheckSubs,
			},
			&cli.StringFlag{
				Name:        "internal-owner",
				Usage:       "git owner/organization to check sub module vulnerabilities",
				Value:       params.INTERNAL_OWNER(),
				Destination: &internalOwner,
			},
		},
		Action: run,
	}
}

func run(ctx context.Context, c *cli.Command) error {
	params.OUTPUT_TYPE = outputType
	if !pathtools.DirExists(projectPath) {
		return fmt.Errorf("project path not found: %s", projectPath)
	}

	store, err := store.New(storePath, internalOwner)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	chk, err := check.New(projectPath, store, internalOwner)
	if err != nil {
		return fmt.Errorf("checker error: %w", err)
	}
	fmt.Println(params.AppName)

	err = chk.Run(check.ProjectCheck)

	if err != nil {
		return err
	}
	if !dontCheckSubs {
		chk.CheckSubs()
	}
	if chk.HasVulnerabilities() && !justWarn {
		return fmt.Errorf("project %s has vulnerabilities", chk.PackageName())
	}
	return nil
}
