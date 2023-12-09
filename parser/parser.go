package parser

import (
	"os"
	"path/filepath"

	"github.com/bensooraj/griffon/blocks"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func ParseHCL(filename string, src []byte, ctx *hcl.EvalContext) (*blocks.Config, error) {
	config := blocks.Config{}
	if err := hclsimple.Decode(filename, src, ctx, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func GetEvalContext() *hcl.EvalContext {
	vars := make(map[string]cty.Value)

	// Region variables
	vars["AMS"] = cty.StringVal("ams")

	// Environment variables
	vars["env"] = cty.ObjectVal(map[string]cty.Value{
		"VULTR_API_KEY": cty.StringVal(os.Getenv("VULTR_API_KEY")),
	})

	// Initialize the data block
	vars["data"] = cty.ObjectVal(map[string]cty.Value{
		string(blocks.RegionBlockType): cty.ObjectVal(map[string]cty.Value{}),
		string(blocks.PlanBlockType):   cty.ObjectVal(map[string]cty.Value{}),
		string(blocks.OSBlockType):     cty.ObjectVal(map[string]cty.Value{}),
	})
	vars[string(blocks.SSHKeyBlockType)] = cty.ObjectVal(map[string]cty.Value{})
	vars[string(blocks.StartupScriptBlockType)] = cty.ObjectVal(map[string]cty.Value{})
	vars[string(blocks.InstanceBlockType)] = cty.ObjectVal(map[string]cty.Value{})

	functions := make(map[string]function.Function)

	// Built-in functions
	functions["uppercase"] = stdlib.UpperFunc // Returns the given string with all Unicode letters translated to their uppercase equivalents
	functions["lowercase"] = stdlib.LowerFunc // Returns the given string with all Unicode letters translated to their lowercase equivalents

	// custom function
	functions["file"] = function.New(&function.Spec{
		Description: "Reads the contents of a file and returns it as a string.",
		Params: []function.Parameter{
			{Type: cty.String},
		},
		Type: func(args []cty.Value) (cty.Type, error) { // or function.StaticReturnType(cty.String),
			return cty.String, nil
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			filename := args[0].AsString()
			fileBuffer, err := os.ReadFile(filepath.Clean(filename))
			if err != nil {
				return cty.StringVal(""), err
			}
			return cty.StringVal(string(fileBuffer)), nil
		},
	})

	return &hcl.EvalContext{
		Variables: vars,
		Functions: functions,
	}
}
