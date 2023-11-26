package blocks

import (
	"errors"

	"github.com/bensooraj/griffon/bodyschema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type SSHKeyBlock struct {
	SSHKey string `hcl:"ssh_key"`
	ResourceBlock
}

var _ Block = (*SSHKeyBlock)(nil)

func (s *SSHKeyBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(bodyschema.DependsOnSchema)
	if diags.HasErrors() {
		return diags
	}
	s.Config = remain

	if attr, ok := content.Attributes["depends_on"]; ok {
		s.DependsOn, diags = ExprAsStringSlice(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}
	return nil
}

func (s *SSHKeyBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := s.Config.PartialContent(bodyschema.SSHKeyBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("ssh_key block must have attributes")
	}

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "ssh_key":
			s.SSHKey = value.AsString()
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func (s *SSHKeyBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
