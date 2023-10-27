package main

import "github.com/hashicorp/hcl/v2/hclsimple"

func ParseHCL(filename string, src []byte) (*Config, error) {
	var config Config
	if err := hclsimple.Decode(filename, src, nil, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
