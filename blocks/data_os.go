package blocks

import (
	"context"
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type OSDataBlock struct {
	VultrID int    `json:"id"`
	OSName  string `json:"name"`
	Arch    string `json:"arch"`
	Family  string `json:"family"`
	DataBlock
}

var _ Block = (*OSDataBlock)(nil)

func (o *OSDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, diags := o.Config.Content(schema.OSFilterSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("os filter block must have attributes")
	}

	var osf schema.OSFilterBlock
	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}
		switch attrName {
		case "type":
			osf.Type = value.AsString()
		case "name":
			osf.Name = value.AsString()
		case "arch":
			osf.Arch = value.AsString()
		case "family":
			osf.Family = value.AsString()
		}
	}
	fmt.Printf("filter: %+v\n", osf)
	return nil
}

// Get
func (o *OSDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
