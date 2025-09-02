package commands

import "github.com/urfave/cli/v3"

func storeCommand() *cli.Command {
	return &cli.Command{
		Name:  "store",
		Usage: "analyses the project throught govulncheck",
	}
}
