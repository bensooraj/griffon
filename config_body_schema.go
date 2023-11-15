package main

import (
	"github.com/hashicorp/hcl/v2"
)

var GriffonBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "region", Required: true},
		{Name: "vultr_api_key", Required: true},
	},
}

var DataBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "filter"},
	},
	Attributes: []hcl.AttributeSchema{},
}

var SSHKeyBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "ssh_key", Required: true},
	},
}

var StartupScriptBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "script", Required: true},
	},
}

var InstanceBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "region", Required: true},
		{Name: "plan", Required: true},
		{Name: "os_id", Required: true},
		{Name: "sshkey_id", Required: true},
		{Name: "script_id", Required: true},
		{Name: "hostname", Required: true},
	},
}

var ConfigSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "griffon", LabelNames: []string{}},
		{Type: "data", LabelNames: []string{"type", "name"}},
		{Type: "ssh_key", LabelNames: []string{"name"}},
		{Type: "startup_script", LabelNames: []string{"name"}},
		{Type: "instance", LabelNames: []string{"name"}},
	},
	Attributes: []hcl.AttributeSchema{},
}

// //////////////////////////////////////////////////////////////////////
var PlanFilterSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "region", Required: true},
		{Name: "vcpu_count", Required: false},
		{Name: "ram", Required: false},
		{Name: "disk", Required: false},
	},
}

type PlanFilterBlock struct {
	Region    string
	VCPUCount int64
	RAM       int64
	Disk      int64
}

var OSFilterSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "type", Required: true},
		{Name: "name", Required: true},
		{Name: "arch", Required: true},
		{Name: "family", Required: true},
	},
}

type OSFilterBlock struct {
	Type   string
	Name   string
	Arch   string
	Family string
}
