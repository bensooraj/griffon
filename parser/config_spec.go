package main

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// ConfigSpec is top-level object spec for the config file.
var ConfigSpec hcldec.ObjectSpec = hcldec.ObjectSpec{
	"griffon": &GriffonSpec,
	"ssh_key": &SSHKeySpec,
}

// GriffonSpec is the spec for the griffon block.
var GriffonSpec hcldec.BlockSpec = hcldec.BlockSpec{
	TypeName: "griffon",
	Nested: &hcldec.ObjectSpec{
		"region": &hcldec.AttrSpec{
			Name:     "region",
			Type:     cty.String,
			Required: true,
		},
		"vultr_api_key": &hcldec.AttrSpec{
			Name:     "vultr_api_key",
			Type:     cty.String,
			Required: true,
		},
	},
}

// SSHKeySpec is the spec for the ssh_key block.
var SSHKeySpec hcldec.BlockListSpec = hcldec.BlockListSpec{
	TypeName: "ssh_key",
	Nested: &hcldec.ObjectSpec{
		"name": &hcldec.BlockLabelSpec{
			Index: 0,
			Name:  "name",
		},
		"ssh_key": &hcldec.AttrSpec{
			Name:     "ssh_key",
			Type:     cty.String,
			Required: true,
		},
	},
}
