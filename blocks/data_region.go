package blocks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type RegionDataBlock struct {
	VID       string   `json:"id" cty:"id"`
	City      string   `json:"city" cty:"city"`
	Country   string   `json:"country" cty:"country"`
	Continent string   `json:"continent" cty:"continent"`
	Options   []string `json:"options" `
	DataBlock
}

var _ Block = (*RegionDataBlock)(nil)

// DataBlock -> RegionDataBlock
func (r *RegionDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	if ctx.Variables["region"].IsNull() {
		return fmt.Errorf("region is not set")
	}
	// r.VultrID = ctx.Variables["region"].AsString()
	return nil
}

// GET
func (r *RegionDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error {
	regions, _, _, err := vc.Region.List(ctx, &govultr.ListOptions{PerPage: 100})
	if err != nil {
		return err
	}
	slog.Debug("Regions", slog.String("block_type", string(r.BlockType())), slog.String("block_name", string(r.BlockName())), slog.Any("regions", regions))
	regionID := evalCtx.Variables["region"].AsString()

	for _, region := range regions {
		if region.ID == regionID {
			r.VID = region.ID
			r.City = region.City
			r.Country = region.Country
			r.Continent = region.Continent
			r.Options = region.Options
			break
		}
	}

	return nil
}

// ToCtyValue
func (r *RegionDataBlock) ToCtyValue() (cty.Value, error) {
	return gocty.ToCtyValue(r, cty.Object(map[string]cty.Type{
		"id":        cty.String,
		"city":      cty.String,
		"country":   cty.String,
		"continent": cty.String,
	}))
}
