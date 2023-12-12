package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bensooraj/griffon/blocks"
	"github.com/bensooraj/griffon/parser"
	"github.com/hashicorp/hcl/v2"
	"github.com/urfave/cli/v2"
	"github.com/vultr/govultr/v3"
)

func main() {
	app := &cli.App{
		Name:  "griffon",
		Usage: "A PoC Vultr cloud automation tool",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Create a set of resources using an HCL file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Required:    true,
						Value:       "griffon.hcl",
						DefaultText: "griffon.hcl",
						Usage:       "HCL file to use",
						Action: func(ctx *cli.Context, s string) error {
							fileInfo, err := os.Stat(s)
							switch {
							case os.IsNotExist(err):
								return cli.Exit(fmt.Sprintf("File %s doesn't exist", s), 1)
							case err != nil:
								return cli.Exit(fmt.Sprintf("Error reading file: %v", err), 1)
							case fileInfo.IsDir():
								return cli.Exit(fmt.Sprintf("%s is a directory", s), 1)
							}
							return nil
						},
					},
				},
				Action: CreateCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// CreateCommand
func CreateCommand(c *cli.Context) error {
	log.Println("create")
	filename := c.String("file")

	// 1. Read the HCL file
	hclFile, err := os.ReadFile(filename)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error reading file: %v", err), 1)
	}
	log.Println(string(hclFile))

	// 2. Parse the HCL file and load the
	var (
		config  *blocks.Config
		evalCtx *hcl.EvalContext = parser.GetEvalContext()
		vc      *govultr.Client
	)

	config, err = parser.ParseWithBodySchema(filename, hclFile, evalCtx)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error parsing file and loading config: %v", err), 1)
	}

	// 3. Evaluate the config
	err = parser.EvaluateConfig(evalCtx, config, vc)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error evaluating config: %v", err), 1)
	}

	return nil
}
