package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type StartupScriptBlock struct {
	GraphID   int64 `json:"graph_id"`
	Config    hcl.Body
	DependsOn []string `json:"depends_on"`
	Name      string   `hcl:"name,label" json:"name"`
	Script    string   `hcl:"script" json:"script"`

	VID           string `json:"id"`
	VDateCreated  string `json:"date_created"`
	VDateModified string `json:"date_modified"`
	VType         string `json:"type"`
}

var _ Block = (*StartupScriptBlock)(nil)

func (s *StartupScriptBlock) ID() int64 {
	return s.GraphID
}

func (s *StartupScriptBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("startup_script block must have attributes")
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

func (s *StartupScriptBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := s.Config.PartialContent(StartupScriptBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("startup_script block must have attributes")
	}

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "script":
			s.Script = value.AsString()
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func (s *StartupScriptBlock) Dependencies() []string {
	return s.DependsOn
}

func (s *StartupScriptBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	fmt.Println("Creating startup script", s.Name)
	ss, _, err := vc.StartupScript.Create(context.Background(), &govultr.StartupScriptReq{
		Name:   s.Name,
		Script: s.Script,
	})
	if err != nil {
		return err
	}

	s.VID = ss.ID
	s.VDateCreated = ss.DateCreated
	s.VDateModified = ss.DateModified
	s.VType = ss.Type

	return nil
}
