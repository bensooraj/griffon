package main

import "github.com/hashicorp/hcl/v2"

var GriffonBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "region", Required: true},
		{Name: "vultr_api_key", Required: true},
	},
}

var SSHKeyBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{Name: "ssh_key", Required: true},
	},
}

var ConfigSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "griffon", LabelNames: []string{}},
		{Type: "ssh_key", LabelNames: []string{"name"}},
	},
	Attributes: []hcl.AttributeSchema{},
}
