package parser

import (
	"fmt"

	"github.com/bensooraj/griffon/blocks"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func ParseHCLUsingSpec(filename string, src []byte, ctx *hcl.EvalContext) (*blocks.Config, error) {
	var config blocks.Config
	var diags hcl.Diagnostics

	parser := hclparse.NewParser()
	// Parse the HCL file
	file, diags := parser.ParseHCL(src, filename)
	if diags.HasErrors() {
		return nil, diags
	}

	val, diags := hcldec.Decode(file.Body, ConfigSpec, ctx)
	if diags.HasErrors() {
		return nil, diags
	}

	config.Data = make(map[blocks.BlockType]map[string]blocks.Block)
	config.Resources = make(map[blocks.BlockType]map[string]blocks.Block)
	for vName, v := range val.AsValueMap() {
		switch vName {
		case "griffon":
			config.Griffon = blocks.GriffonBlock{
				Region:      v.GetAttr("region").AsString(),
				VultrAPIKey: v.GetAttr("vultr_api_key").AsString(),
			}
		case "ssh_key":
			sshKeys := v.AsValueSlice()
			for _, sshKey := range sshKeys {
				r := &blocks.SSHKeyBlock{SSHKey: sshKey.GetAttr("ssh_key").AsString()}
				r.Name = sshKey.GetAttr("name").AsString()
				r.Type = "ssh_key"
				config.AddResource(r)
			}
		default:
			return nil, fmt.Errorf("unknown block type %q", vName)
		}
	}

	return &config, nil
}
