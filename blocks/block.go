package blocks

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

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
	Create(ctx *hcl.EvalContext, vc *govultr.Client) error
}
