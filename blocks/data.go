package blocks

import (
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/bodyschema"
	"github.com/hashicorp/hcl/v2"
)

type DataBlock struct {
	GraphID   int64
	Type      string `hcl:"type,label"`
	Name      string `hcl:"name,label"`
	Config    hcl.Body
	DependsOn []string
}

func (d *DataBlock) ID() int64 {
	return d.GraphID
}

func (d *DataBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(bodyschema.DependsOnSchema)
	if diags.HasErrors() {
		return diags
	}
	if attr, ok := content.Attributes["depends_on"]; ok {
		d.DependsOn, diags = ExprAsStringSlice(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}

	filterBodyContent, _, diags := remain.PartialContent(bodyschema.DataBlockSchema)
	if diags.HasErrors() {
		return diags
	}
	switch d.Type {
	case "region":
		// do nothing
	case "plan", "os":
		if len(filterBodyContent.Blocks) != 1 {
			return fmt.Errorf("%s block %s must have one filter block", d.Type, d.Name)
		}
		d.Config = filterBodyContent.Blocks[0].Body
	default:
		return errors.New("unknown data type " + d.Type + " with name " + d.Name)
	}
	return nil
}

func (d *DataBlock) Dependencies() []string {
	return d.DependsOn
}