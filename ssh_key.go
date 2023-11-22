package main

import (
	"errors"

	"github.com/hashicorp/hcl/v2"
)

type SSHKeyBlock struct {
	GraphID   int64
	Name      string `hcl:"name,label"`
	SSHKey    string `hcl:"ssh_key"`
	Config    hcl.Body
	DependsOn []string
}

func (s *SSHKeyBlock) ID() int64 {
	return s.GraphID
}
func (s *SSHKeyBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
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
	content, _, diags := s.Config.PartialContent(SSHKeyBlockSchema)
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

func (s *SSHKeyBlock) Dependencies() []string {
	return s.DependsOn
}
