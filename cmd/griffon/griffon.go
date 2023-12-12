package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
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
	hclFile, err := os.ReadFile(c.String("file"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Error reading file: %v", err), 1)
	}
	log.Println(string(hclFile))
	return nil
}
