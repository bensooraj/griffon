package main

import "github.com/hashicorp/hcl/v2"

type Config struct {
	Griffon        GriffonBlock                  `hcl:"griffon,block"`
	SSHKeys        map[string]SSHKeyBlock        `hcl:"ssh_key,block"`
	StartupScripts map[string]StartupScriptBlock `hcl:"startup_script,block"`
	Instances      map[string]InstanceBlock      `hcl:"instance,block"`
	DataBlocks     map[string]map[string]Block   `hcl:"data,block"`
}

type Block interface {
	ID() int64
	// Separate the block into its configuration and dependencies
	PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error
	// Process the configuration
	ProcessConfiguration(ctx *hcl.EvalContext) error
	// Get Dependencies
	Dependencies() []string
	// Execute the block by making API calls
	// Execute(ctx *hcl.EvalContext) error
}
