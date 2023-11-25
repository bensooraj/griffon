package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
)

type InstanceBlock struct {
	GraphID         int64
	Name            string            `hcl:"name,label"`
	Region          string            `hcl:"region,attr"`
	Plan            string            `hcl:"plan,attr"`
	OS              string            `hcl:"os,attr"`
	SshKeyID        string            `hcl:"ssh_key_id,attr"`
	StartupScriptID string            `hcl:"startup_script_id,attr"`
	Hostname        string            `hcl:"hostname,attr"`
	Tag             map[string]string `hcl:"tag,attr"`
	Config          hcl.Body
	DependsOn       []string
}

var _ Block = (*InstanceBlock)(nil)

func (i *InstanceBlock) ID() int64 {
	return i.GraphID
}

func (i *InstanceBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
	if diags.HasErrors() {
		return diags
	}
	i.Config = remain

	if attr, ok := content.Attributes["depends_on"]; ok {
		i.DependsOn, diags = ExprAsStringSlice(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}
	return nil
}

func (i *InstanceBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := i.Config.PartialContent(InstanceBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("instance block must have attributes")
	}

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "region":
			i.Region = value.AsString()
		case "plan":
			i.Plan = value.AsString()
		case "os":
			i.OS = value.AsString()
		case "ssh_key_id":
			i.SshKeyID = value.AsString()
		case "startup_script_id":

		case "hostname":
			i.Hostname = value.AsString()
		case "tag":
			i.Tag = make(map[string]string)

			fmt.Println("tag:", value.AsString(), value.AsValueMap())
			for key, ctyVal := range value.AsValueMap() {
				i.Tag[key] = ctyVal.AsString()
			}
		}
	}

	return nil
}

func (i *InstanceBlock) Dependencies() []string {
	return i.DependsOn
}

func (i *InstanceBlock) Create(ctx *hcl.EvalContext, vc *govultr.Client) error {
	return nil
}
