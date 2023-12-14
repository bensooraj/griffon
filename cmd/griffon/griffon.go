package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/bensooraj/griffon/blocks"
	"github.com/bensooraj/griffon/parser"
	"github.com/hashicorp/hcl/v2"
	"github.com/urfave/cli/v2"
	"github.com/vultr/govultr/v3"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}

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
				OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// CreateCommand
func CreateCommand(c *cli.Context) error {
	filename := c.String("file")
	// 1. Read the HCL file
	hclFile, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("[CMD CREATE] Error reading file", slog.String("filename", filename), slog.String("error", err.Error()))
		return cli.Exit(fmt.Sprintf("Error reading file: %v", err), 1)
	}
	slog.Debug("[CMD CREATE] Read HCL file", slog.String("filename", filename), slog.String("file", string(hclFile)))

	// 2. Parse the HCL file and load the
	var (
		config  *blocks.Config
		evalCtx *hcl.EvalContext = parser.GetEvalContext()
		vc      *govultr.Client
	)

	config, err = parser.ParseWithBodySchema(filename, hclFile, evalCtx)
	if err != nil {
		slog.Error("[CMD CREATE] Error parsing file and loading config", slog.String("filename", filename), slog.String("error", err.Error()))
		return cli.Exit(fmt.Sprintf("Error parsing file and loading config: %v", err), 1)
	}

	// 3. Evaluate the config
	err = parser.EvaluateConfig(evalCtx, config, vc)
	if err != nil {
		slog.Error("[CMD CREATE] Error evaluating config", slog.String("filename", filename), slog.String("error", err.Error()))
		return cli.Exit(fmt.Sprintf("Error evaluating config: %v", err), 1)
	}

	return nil
}
