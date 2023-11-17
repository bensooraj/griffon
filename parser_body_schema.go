package main

import (
	"errors"
	"fmt"
	"log"

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
	// Calculate dependencies and store them in s.Config
	err := s.CalculateDependency(block, ctx)
	if err != nil {
		return err
	}

	content, _, diags := block.Body.PartialContent(StartupScriptBlockSchema)
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

func (s *StartupScriptBlock) CalculateDependency(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("startup_script block must have attributes")
	}
	s.Config = remain

	if attr, ok := content.Attributes["depends_on"]; ok {
		s.DependsOn, diags = ExprAsMap(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}
	return nil
}

func (d *DataBlock) FromHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, diags := block.Body.Content(DataBlockSchema)
	if diags.HasErrors() {
		return diags
	}

	// Labels
	d.Type = block.Labels[0]
	d.Name = block.Labels[1]
	// Blocks
	var filterBlock *hcl.Block
	blocks := content.Blocks.OfType("filter")
	if d.Type != "region" && len(blocks) != 1 {
		return errors.New("data block must have one filter block")
	} else if len(blocks) == 1 {
		filterBlock = blocks[0]
	}

	switch d.Type {
	case "region":
		fmt.Println("- data::region -")
	case "plan":
		fmt.Println("- data::plan -")
		filterContent, diags := filterBlock.Body.Content(PlanFilterSchema)
		if diags.HasErrors() {
			return diags
		}
		var pf PlanFilterBlock
		for attrName, attr := range filterContent.Attributes {
			value, diags := attr.Expr.Value(ctx)
			if diags.HasErrors() {
				return diags
			}
			switch attrName {
			case "region":
				pf.Region = value.AsString()
			case "vcpu_count":
				pf.VCPUCount, _ = value.AsBigFloat().Int64()
			case "ram":
				pf.RAM, _ = value.AsBigFloat().Int64()
			case "disk":
				pf.Disk, _ = value.AsBigFloat().Int64()
			}
		}
		fmt.Printf("filter: %+v\n", pf)
	case "os":
		fmt.Println("- data::os -")
		filterContent, diags := filterBlock.Body.Content(OSFilterSchema)
		if diags.HasErrors() {
			return diags
		}
		var osf OSFilterBlock
		for attrName, attr := range filterContent.Attributes {
			value, diags := attr.Expr.Value(ctx)
			if diags.HasErrors() {
				return diags
			}
			switch attrName {
			case "type":
				osf.Type = value.AsString()
			case "name":
				osf.Name = value.AsString()
			case "arch":
				osf.Arch = value.AsString()
			case "family":
				osf.Family = value.AsString()
			}
		}
		fmt.Printf("filter: %+v\n", osf)
	default:
		return errors.New("unknown data type " + d.Type)
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
		fmt.Println("blockName:", blockName)
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
		case "data":
			for _, hclBlock := range hclBlocks {
				var data DataBlock
				if err := data.FromHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
			}
		default:
			fmt.Println("unknown block type", blockName)
		}
	}

	fmt.Println()
	log.Println("~~~~~~~~~~~ section 2 ~~~~~~~~~~~")
	//
	for blockName, hclBlocks := range blocks {
		switch blockName {
		case "startup_script":
			for _, hclBlock := range hclBlocks {
				var startupScript StartupScriptBlock
				if err := startupScript.CalculateDependency(hclBlock, ctx); err != nil {
					return nil, err
				}
			}
		default:
			fmt.Println("unknown block type", blockName)
		}
	}
	fmt.Println()

	return &config, nil
}
