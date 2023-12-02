package blocks

import (
	"context"
	"fmt"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type ResourceBlock struct {
	GraphID   int64     `json:"graph_id"`
	Type      BlockType `hcl:"type,label" json:"type"`
	Name      string    `hcl:"name,label" json:"name"`
	Config    hcl.Body
	DependsOn []string `json:"depends_on"`
}

func (r *ResourceBlock) ID() int64 {
	return r.GraphID
}

func (r *ResourceBlock) BlockType() BlockType {
	return r.Type
}

func (r *ResourceBlock) BlockName() string {
	return r.Name
}

func (r *ResourceBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(schema.DependsOnSchema)
	if diags.HasErrors() {
		return diags
	}
	r.Config = remain

	if attr, ok := content.Attributes["depends_on"]; ok {
		r.DependsOn, diags = ExprAsStringSlice(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}
	fmt.Println(block.Type, r.Name, " depends_on:", r.DependsOn)

	r.Type = BlockType(block.Type)

	return nil
}

func (r *ResourceBlock) Dependencies() []string {
	return r.DependsOn
}

// Implement Get
func (r *ResourceBlock) Get(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) (*hcl.EvalContext, error) {
	return nil, nil
}
