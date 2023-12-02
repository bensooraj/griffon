package blocks

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
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
	if ctx.Variables["region"].IsNull() {
		return fmt.Errorf("region is not set")
	}
	// r.VultrID = ctx.Variables["region"].AsString()
	return nil
}

// GET
func (r *RegionDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) (*hcl.EvalContext, error) {
	regions, meta, _, err := vc.Region.List(ctx, &govultr.ListOptions{PerPage: 100})
	if err != nil {
		return nil, err
	}
	fmt.Println("RegionDataBlock::Get::meta:", meta)
	regionID := evalCtx.Variables["region"].AsString()

	for _, region := range regions {
		if region.ID == regionID {
			r.VultrID = region.ID
			r.City = region.City
			r.Country = region.Country
			r.Continent = region.Continent
			r.Options = region.Options
			break
		}
	}

	// Update the evaluation context variables
	evalCtx.Variables = map[string]cty.Value{
		"data": cty.ObjectVal(map[string]cty.Value{
			"region": cty.ObjectVal(map[string]cty.Value{
				r.Name: cty.ObjectVal(map[string]cty.Value{
					"id":        cty.StringVal(r.VultrID),
					"city":      cty.StringVal(r.City),
					"country":   cty.StringVal(r.Country),
					"continent": cty.StringVal(r.Continent),
				}),
			}),
		}),
	}

	fmt.Printf("....(data.region.%s) Evaluation context: %s\n", r.Name, evalCtx.Variables["data"].GoString())

	return nil, nil
}
