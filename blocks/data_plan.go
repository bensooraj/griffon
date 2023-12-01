package blocks

import (
	"context"
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
)

type PlanDataBlock struct {
	VultrID     string   `json:"id"`
	VCPUCount   int      `json:"vcpu_count"`
	RAM         int      `json:"ram"`
	Disk        int      `json:"disk"`
	DiskCount   int      `json:"disk_count"`
	Bandwidth   int      `json:"bandwidth"`
	MonthlyCost float32  `json:"monthly_cost"`
	PlanType    string   `json:"type"`
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
func (p *PlanDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error {
	plans, meta, _, err := vc.Plan.List(ctx, p.filter.Type, &govultr.ListOptions{PerPage: 100})
	if err != nil {
		return err
	}
	fmt.Println("PlanDataBlock::Get::meta:", meta)
	found := false

	for _, plan := range plans {
		if int64(plan.VCPUCount) == p.filter.VCPUCount &&
			int64(plan.RAM) == p.filter.RAM &&
			int64(plan.Disk) == p.filter.Disk {

			for _, loc := range plan.Locations {
				if loc == p.filter.Region {
					p.VultrID = plan.ID
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
		return ErrorDataNotFound
	}

	// Update the evaluation context variables
	newEvalCtx := evalCtx.NewChild()
	newEvalCtx.Variables = map[string]cty.Value{
		"data": cty.ObjectVal(map[string]cty.Value{
			"plan": cty.ObjectVal(map[string]cty.Value{
				p.Name: cty.ObjectVal(map[string]cty.Value{
					"id":           cty.StringVal(p.VultrID),
					"vcpu_count":   cty.NumberIntVal(int64(p.VCPUCount)),
					"ram":          cty.NumberIntVal(int64(p.RAM)),
					"disk":         cty.NumberIntVal(int64(p.Disk)),
					"disk_count":   cty.NumberIntVal(int64(p.DiskCount)),
					"bandwidth":    cty.NumberIntVal(int64(p.Bandwidth)),
					"monthly_cost": cty.NumberFloatVal(float64(p.MonthlyCost)),
					"type":         cty.StringVal(p.PlanType),
				}),
			}),
		}),
	}

	fmt.Printf("....(data.plan.%s) Evaluation context: %s\n", p.Name, evalCtx.Variables["data"].GoString())

	return nil
}
