package blocks

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type RegionDataBlock struct {
	VultrID   string   `json:"id"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Continent string   `json:"continent"`
	Options   []string `json:"options"`
	DataBlock
}

var _ Block = (*RegionDataBlock)(nil)

// DataBlock -> RegionDataBlock
func (r *RegionDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	// Nothing to do here really
	return nil
}

func (r *RegionDataBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
