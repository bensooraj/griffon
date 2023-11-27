package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func ParseHCLUsingSpec(filename string, src []byte, ctx *hcl.EvalContext) (*Config, error) {
	var config Config
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

	for vName, v := range val.AsValueMap() {
		switch vName {
		case "griffon":
			config.Griffon = GriffonBlock{
				Region:      v.GetAttr("region").AsString(),
				VultrAPIKey: v.GetAttr("vultr_api_key").AsString(),
			}
		case "ssh_key":
			sshKeys := v.AsValueSlice()
			for _, sshKey := range sshKeys {
				config.SSHKeys[sshKey.GetAttr("name").AsString()] = SSHKeyBlock{
					Name:   sshKey.GetAttr("name").AsString(),
					SSHKey: sshKey.GetAttr("ssh_key").AsString(),
				}
			}
		default:
			return nil, fmt.Errorf("unknown block type %q", vName)
		}
	}

	return &config, nil
}
