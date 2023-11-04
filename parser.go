package main

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func ParseHCL(filename string, src []byte, ctx *hcl.EvalContext) (*Config, error) {
	config := Config{}
	if err := hclsimple.Decode(filename, src, ctx, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getEvalContext() *hcl.EvalContext {
	vars := make(map[string]cty.Value)

	// Region variables
	vars["AMS"] = cty.StringVal("ams")

	// Environment variables
	vars["env"] = cty.ObjectVal(map[string]cty.Value{
		"VULTR_API_KEY": cty.StringVal(os.Getenv("VULTR_API_KEY")),
	})

	functions := make(map[string]function.Function)

	// Built-in functions
	functions["uppercase"] = stdlib.UpperFunc // Returns the given string with all Unicode letters translated to their uppercase equivalents
	functions["lowercase"] = stdlib.LowerFunc // Returns the given string with all Unicode letters translated to their lowercase equivalents

	return &hcl.EvalContext{
		Variables: vars,
		Functions: functions,
	}
}
