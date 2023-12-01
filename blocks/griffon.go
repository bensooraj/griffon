package blocks

import (
	"errors"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type GriffonBlock struct {
	GraphID     int64
	Region      string `hcl:"region,attr"`
	VultrAPIKey string `hcl:"vultr_api_key"`
	Block
}

var _ Block = (*GriffonBlock)(nil)

func (g *GriffonBlock) ID() int64 {
	return g.GraphID
}

func (g *GriffonBlock) BlockType() BlockType {
	return GriffonBlockType
}

func (g *GriffonBlock) BlockName() string {
	return "griffon"
}

func (g *GriffonBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, diags := block.Body.Content(schema.GriffonBlockSchema)
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
	ctx.Variables["region"] = cty.StringVal(g.Region)
	ctx.Variables["vultr_api_key"] = cty.StringVal(g.VultrAPIKey)
	return nil
}

func (g *GriffonBlock) Dependencies() []string {
	return nil
}
