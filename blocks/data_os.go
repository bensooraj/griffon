package blocks

import (
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/bodyschema"
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

func (o *OSDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, diags := o.Config.Content(bodyschema.OSFilterSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("os filter block must have attributes")
	}

	var osf bodyschema.OSFilterBlock
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

func (o *OSDataBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
