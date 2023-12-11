package blocks

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
)

type Block interface {
	ID() int64
	BlockType() BlockType
	BlockName() string
	// Separate the block into its configuration and dependencies
	PreProcessHCLBlock(block *hcl.Block, evalCtx *hcl.EvalContext) error
	// Process the configuration
	ProcessConfiguration(evalCtx *hcl.EvalContext) error
	// Get Dependencies
	Dependencies() []string
	// Get the data block
	Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error
	// Create the block by making API calls
	Create(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error
	// ToCtyValue
	ToCtyValue() (cty.Value, error)
}
