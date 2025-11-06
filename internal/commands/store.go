package commands

import (
	"context"
	"fmt"

	"github.com/guionardo/govuln/internal/store"
	"github.com/urfave/cli/v3"
)

func storeCommand() *cli.Command {
	return &cli.Command{
		Name:  "store",
		Usage: "data analysed",
		Commands: []*cli.Command{
			storeShowCommand(),
		},
	}
}

func storeShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "display analysed data",
		Action: func(ctx context.Context, c *cli.Command) error {
			s, err := store.New(storePath, internalOwner)
			if err != nil {
				return fmt.Errorf("store error: %w", err)
			}
			s.ShowInfo()
			return nil
		},
	}
}
