package blocks

import (
	"errors"

	"github.com/bensooraj/griffon/bodyschema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type GriffonBlock struct {
	GraphID     int64
	Region      string `hcl:"region,attr"`
	VultrAPIKey string `hcl:"vultr_api_key"`
}

var _ Block = (*GriffonBlock)(nil)

func (g *GriffonBlock) ID() int64 {
	return g.GraphID
}

func (g *GriffonBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, diags := block.Body.Content(bodyschema.GriffonBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("griffon block must have attributes")
	}

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "region":
			g.Region = value.AsString()
		case "vultr_api_key":
			g.VultrAPIKey = value.AsString()
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func (g *GriffonBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	// Nothing to do here really
	return nil
}

func (g *GriffonBlock) Dependencies() []string {
	return nil
}

func (g *GriffonBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
