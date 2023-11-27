package main

import "github.com/bensooraj/griffon/blocks"

type Config struct {
	Griffon   blocks.GriffonBlock             `hcl:"griffon,block" json:"griffon"`
	Data      map[string]blocks.DataBlock     `hcl:"data,block" json:"data"`
	Resources map[string]blocks.ResourceBlock `hcl:"resources,block" json:"resources"`
}
