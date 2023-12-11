package blocks

import (
	"context"
	"errors"
	"strings"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type OSDataBlock struct {
	VID    int    `json:"id" cty:"id"`
	OSName string `json:"name" cty:"name"`
	Arch   string `json:"arch" cty:"arch"`
	Family string `json:"family" cty:"family"`
	filter schema.OSFilterBlock
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
	o.filter = osf
	return nil
}

// Get
func (o *OSDataBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) (*hcl.EvalContext, error) {
	oss, _, _, err := vc.OS.List(ctx, &govultr.ListOptions{PerPage: 100})
	if err != nil {
		return nil, err
	}

	found := false

	for _, os := range oss {
		if os.Family == o.filter.Family &&
			os.Arch == o.filter.Arch &&
			strings.Contains(os.Name, o.filter.Name) {
			found = true
			o.OSName = os.Name
			o.VID = os.ID
			o.Arch = os.Arch
			o.Family = os.Family

			found = true
			break
		}
	}
	if !found {
		return nil, ErrorDataNotFound
	}

	return nil, nil
}

// ToCtyValue
func (o *OSDataBlock) ToCtyValue() (cty.Value, error) {
	return gocty.ToCtyValue(o, cty.Object(map[string]cty.Type{
		"id":     cty.Number,
		"name":   cty.String,
		"arch":   cty.String,
		"family": cty.String,
	}))
}
