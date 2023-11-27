package blocks

import (
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type PlanDataBlock struct {
	VultrID     string   `json:"id"`
	VcpuCount   int      `json:"vcpu_count"`
	Ram         int      `json:"ram"`
	Disk        int      `json:"disk"`
	DiskCount   int      `json:"disk_count"`
	Bandwidth   int      `json:"bandwidth"`
	MonthlyCost int      `json:"monthly_cost"`
	PlanType    string   `json:"type"`
	Locations   []string `json:"locations"`
	DataBlock
}

var _ Block = (*PlanDataBlock)(nil)

func (p *PlanDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, diags := p.Config.Content(schema.PlanFilterSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("plan filter block must have attributes")
	}

	var pf schema.PlanFilterBlock
	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}
		switch attrName {
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
	fmt.Printf("filter: %+v\n", pf)
	return nil
}

func (p *PlanDataBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
