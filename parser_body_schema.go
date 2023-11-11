package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func (g *GriffonBlock) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, diags := block.Body.Content(GriffonBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("griffon block must have attributes")
	}

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "region":
			g.Region = value.AsString()
		case "vultr_api_key":
			g.VultrAPIKey = value.AsString()
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func (s *SSHKeyBlock) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, diags := block.Body.Content(SSHKeyBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("ssh_key block must have attributes")
	}

	s.Name = block.Labels[0]

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "ssh_key":
			s.SSHKey = value.AsString()
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func (s *StartupScriptBlock) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, diags := block.Body.Content(StartupScriptBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("startup_script block must have attributes")
	}

	s.Name = block.Labels[0]

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "script":
			s.Script = value.AsString()
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func ParseHCLUsingBodySchema(filename string, src []byte, ctx *hcl.EvalContext) (*Config, error) {
	config := Config{}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	bodyContent, diags := file.Body.Content(ConfigSchema)
	if diags.HasErrors() {
		return nil, diags
	}

	if len(bodyContent.Blocks) == 0 {
		return nil, errors.New("no blocks found")
	}

	blocks := bodyContent.Blocks.ByType()
	for blockName, hclBlocks := range blocks {
		switch blockName {
		case "griffon":
			if len(hclBlocks) != 1 {
				return nil, errors.New("only one griffon block allowed")
			}
			var griffon GriffonBlock
			if err := griffon.FromHCLBlock(hclBlocks[0], ctx); err != nil {
				return nil, err
			}
			config.Griffon = griffon
		case "ssh_key":
			for _, hclBlock := range hclBlocks {
				var sshKey SSHKeyBlock
				if err := sshKey.FromHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				config.SSHKeys = append(config.SSHKeys, sshKey)
			}
		case "startup_script":
			for _, hclBlock := range hclBlocks {
				var startupScript StartupScriptBlock
				if err := startupScript.FromHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				config.StartupScripts = append(config.StartupScripts, startupScript)
			}
		default:
			fmt.Println("unknown block type", blockName)
		}
	}

	return &config, nil
}
