package main

import (
	"errors"

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
				sshKeyContent, diags := hclBlock.Body.Content(SSHKeyBlockSchema)
				switch {
				case diags.HasErrors():
					return nil, diags
				case len(sshKeyContent.Attributes) == 0:
					return nil, errors.New("ssh_key block must have attributes")
				}

				sshKeyBlock := SSHKeyBlock{}
				sshKeyBlock.Name = hclBlock.Labels[0]
				for attrName, attr := range sshKeyContent.Attributes {
					value, diags := attr.Expr.Value(ctx)
					if diags.HasErrors() {
						return nil, diags
					}

					switch attrName {
					case "ssh_key":
						sshKeyBlock.SSHKey = value.AsString()
					default:
						return nil, errors.New("unknown attribute " + attrName)
					}
				}

				config.SSHKeys = append(config.SSHKeys, sshKeyBlock)
			}
		}
	}

	return &config, nil
}
