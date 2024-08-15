package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/grevych/branchswapper"
)

func main() {
	vcs := branchswapper.NewGitExecutor()
	swapper := branchswapper.NewBranchSwapper(vcs)

	app := &cli.App{
		Name:  "branchswap",
		Usage: "Swap git branches",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "list",
				Usage:   "List swapped branches",
				Aliases: []string{"ls"},
			},
			&cli.IntFlag{
				Name:    "index",
				Usage:   "Index",
				Aliases: []string{"i"},
				Value:   -1,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if err := swapper.Load(); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if cCtx.Bool("list") {
				for index, branch := range swapper.GetStack() {
					fmt.Printf("%d: %s\n", index, branch)
				}
				return nil
			}

			if cCtx.Int("index") != -1 {
				if err := swapper.SwapFromStack(cCtx.Int("index")); err != nil {
					return cli.Exit(err.Error(), 1)
				}
				return nil
			}

			branch := cCtx.Args().First()
			if err := swapper.Swap(branch); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
