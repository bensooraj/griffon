package blocks

import (
	"context"
	"fmt"

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

// GET
func (r *RegionDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error {
	region, meta, _, err := vc.Region.List(ctx, &govultr.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Println("RegionDataBlock::Get::meta:", meta)

	for _, r := range region {
		fmt.Println("RegionDataBlock::Get::region:", r)
	}
	return nil
}
