package blocks

import (
	"context"
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type PlanDataBlock struct {
	VID         string   `json:"id" cty:"id"`
	VCPUCount   int      `json:"vcpu_count" cty:"vcpu_count"`
	RAM         int      `json:"ram" cty:"ram"`
	Disk        int      `json:"disk" cty:"disk"`
	DiskCount   int      `json:"disk_count" cty:"disk_count"`
	Bandwidth   int      `json:"bandwidth" cty:"bandwidth"`
	MonthlyCost float32  `json:"monthly_cost" cty:"monthly_cost"`
	PlanType    string   `json:"type" cty:"type"`
	Locations   []string `json:"locations"`
	filter      schema.PlanFilterBlock
	DataBlock
}

var _ Block = (*PlanDataBlock)(nil)

func (p *PlanDataBlock) ProcessConfiguration(evalCtx *hcl.EvalContext) error {
	content, diags := p.Config.Content(schema.PlanFilterSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("plan filter block must have attributes")
	}

	var pf schema.PlanFilterBlock
	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(evalCtx)
		if diags.HasErrors() {
			return diags
		}
		switch attrName {
		case "type":
			pf.Type = value.AsString()
		case "region":
			pf.Region = value.AsString()
		case "vcpu_count":
			pf.VCPUCount, _ = value.AsBigFloat().Int64()
		case "ram":
			pf.RAM, _ = value.AsBigFloat().Int64()
		case "disk":
			pf.Disk, _ = value.AsBigFloat().Int64()
		}
	}
	p.filter = pf
	fmt.Printf("filter: %+v\n", pf)
	return nil
}

// Get
func (p *PlanDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) (*hcl.EvalContext, error) {
	plans, meta, _, err := vc.Plan.List(ctx, p.filter.Type, &govultr.ListOptions{PerPage: 100})
	if err != nil {
		return nil, err
	}
	fmt.Println("PlanDataBlock::Get::meta:", meta)
	found := false

	for _, plan := range plans {
		if int64(plan.VCPUCount) == p.filter.VCPUCount &&
			int64(plan.RAM) == p.filter.RAM &&
			int64(plan.Disk) == p.filter.Disk {

			for _, loc := range plan.Locations {
				if loc == p.filter.Region {
					p.VID = plan.ID
					p.VCPUCount = plan.VCPUCount
					p.RAM = plan.RAM
					p.Disk = plan.Disk
					p.DiskCount = plan.DiskCount
					p.Bandwidth = plan.Bandwidth
					p.MonthlyCost = plan.MonthlyCost
					p.PlanType = plan.Type
					p.Locations = plan.Locations
					found = true
					break
				} else {
					fmt.Println("loc:", loc, "p.filter.Region:", p.filter.Region, plan.Locations)
				}
			}
		}
	}

	if !found {
		return nil, ErrorDataNotFound
	}

	fmt.Printf("\n....(data.plan.%s) Evaluation context: %s\n", p.Name, evalCtx.Variables["data"].GoString())

	return nil, nil
}

// ToCtyValue
func (p *PlanDataBlock) ToCtyValue() (cty.Value, error) {
	return gocty.ToCtyValue(p, cty.Object(map[string]cty.Type{
		"id":           cty.String,
		"vcpu_count":   cty.Number,
		"ram":          cty.Number,
		"disk":         cty.Number,
		"disk_count":   cty.Number,
		"bandwidth":    cty.Number,
		"monthly_cost": cty.Number,
		"type":         cty.String,
	}))
}
