package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// GriffonBlock
func (g *GriffonBlock) ID() int64 {
	return g.GraphID
}

func (g *GriffonBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
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

func (g *GriffonBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	// Nothing to do here really
	return nil
}

func (g *GriffonBlock) Dependencies() map[string][]string {
	return nil
}

// SSHKeyBlock
func (s *SSHKeyBlock) ID() int64 {
	return s.GraphID
}
func (s *SSHKeyBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
	if diags.HasErrors() {
		return diags
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

func (s *SSHKeyBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := s.Config.PartialContent(SSHKeyBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("ssh_key block must have attributes")
	}

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

func (s *SSHKeyBlock) Dependencies() map[string][]string {
	return s.DependsOn
}

// StartupScriptBlock
func (s *StartupScriptBlock) ID() int64 {
	return s.GraphID
}
func (s *StartupScriptBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
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

func (s *StartupScriptBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := s.Config.PartialContent(StartupScriptBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("startup_script block must have attributes")
	}

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

func (s *StartupScriptBlock) Dependencies() map[string][]string {
	return s.DependsOn
}

// IntanceBlock
func (i *InstanceBlock) ID() int64 {
	return i.GraphID
}

func (i *InstanceBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
	if diags.HasErrors() {
		return diags
	}
	i.Config = remain

	if attr, ok := content.Attributes["depends_on"]; ok {
		i.DependsOn, diags = ExprAsMap(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}
	return nil
}

func (i *InstanceBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := i.Config.PartialContent(InstanceBlockSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("instance block must have attributes")
	}

	for attrName, attr := range content.Attributes {
		value, diags := attr.Expr.Value(ctx)
		if diags.HasErrors() {
			return diags
		}

		switch attrName {
		case "region":
			i.Region = value.AsString()
		case "plan":
			i.Plan = value.AsString()
		case "os":
			i.OS = value.AsString()
		case "ssh_key_id":
			i.SshKeyID = value.AsString()
		case "startup_script_id":

		case "hostname":
			i.Hostname = value.AsString()
		case "tag":
			i.Tag = make(map[string]string)

			fmt.Println("tag:", value.AsString(), value.AsValueMap())
			for key, ctyVal := range value.AsValueMap() {
				i.Tag[key] = ctyVal.AsString()
			}
		}
	}

	return nil
}

func (i *InstanceBlock) Dependencies() map[string][]string {
	return i.DependsOn
}

// DataBlock
func (d *DataBlock) ID() int64 {
	return d.GraphID
}

func (d *DataBlock) PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error {
	content, remain, diags := block.Body.PartialContent(DependsOnSchema)
	if diags.HasErrors() {
		return diags
	}
	if attr, ok := content.Attributes["depends_on"]; ok {
		d.DependsOn, diags = ExprAsMap(attr.Expr)
		if diags.HasErrors() {
			return diags
		}
	}

	filterBodyContent, _, diags := remain.PartialContent(DataBlockSchema)
	if diags.HasErrors() {
		return diags
	}
	switch d.Type {
	case "region":
		// do nothing
	case "plan", "os":
		if len(filterBodyContent.Blocks) != 1 {
			return fmt.Errorf("%s block %s must have one filter block", d.Type, d.Name)
		}
		d.Config = filterBodyContent.Blocks[0].Body
	default:
		return errors.New("unknown data type " + d.Type + " with name " + d.Name)
	}
	return nil
}

func (d *DataBlock) Dependencies() map[string][]string {
	return d.DependsOn
}

// DataBlock -> RegionDataBlock
func (r *RegionDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	// Nothing to do here really
	return nil
}

// DataBlock -> PlanDataBlock
func (p *PlanDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, diags := p.Config.Content(PlanFilterSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("plan filter block must have attributes")
	}

	var pf PlanFilterBlock
	for attrName, attr := range content.Attributes {
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
	return nil
}

// DataBlock -> OSDataBlock
func (o *OSDataBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, diags := o.Config.Content(OSFilterSchema)
	switch {
	case diags.HasErrors():
		return diags
	case len(content.Attributes) == 0:
		return errors.New("os filter block must have attributes")
	}

	var osf OSFilterBlock
	for attrName, attr := range content.Attributes {
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
	return nil
}

func ParseHCLUsingBodySchema(filename string, src []byte, ctx *hcl.EvalContext) (*Config, error) {
	config := Config{
		SSHKeys:        make(map[string]SSHKeyBlock),
		StartupScripts: make(map[string]StartupScriptBlock),
		Instances:      make(map[string]InstanceBlock),
		DataBlocks:     make(map[string]map[string]Block),
	}

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

	blockPathToGraphID := make(map[string]int64)

	dependencyGraph := NewGraph()
	var dGraphNodeCount int64 = 0

	blocks := bodyContent.Blocks.ByType()
	for blockName, hclBlocks := range blocks {
		fmt.Println("blockName:", blockName)
		switch blockName {
		case "griffon":
			if len(hclBlocks) != 1 {
				return nil, errors.New("only one griffon block allowed")
			}
			var griffon GriffonBlock
			if err := griffon.PreProcessHCLBlock(hclBlocks[0], ctx); err != nil {
				return nil, err
			}
			griffon.GraphID = 0
			config.Griffon = griffon
			blockPathToGraphID[blockName] = griffon.GraphID

			dependencyGraph.AddNode(&griffon)
		case "ssh_key":
			for _, hclBlock := range hclBlocks {
				var sshKey SSHKeyBlock
				sshKey.Name = hclBlock.Labels[0]

				if err := sshKey.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				sshKey.GraphID = dGraphNodeCount
				blockPathToGraphID[BuildBlockPath(blockName, sshKey.Name)] = sshKey.GraphID

				config.SSHKeys[sshKey.Name] = sshKey
				dependencyGraph.AddNode(&sshKey)
			}
		case "startup_script":
			for _, hclBlock := range hclBlocks {
				var startupScript StartupScriptBlock
				startupScript.Name = hclBlock.Labels[0]

				if err := startupScript.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				startupScript.GraphID = dGraphNodeCount
				blockPathToGraphID[BuildBlockPath(blockName, startupScript.Name)] = startupScript.GraphID

				config.StartupScripts[startupScript.Name] = startupScript
				dependencyGraph.AddNode(&startupScript)
			}
		case "instance":
			for _, hclBlock := range hclBlocks {
				var instance InstanceBlock
				instance.Name = hclBlock.Labels[0]

				if err := instance.PreProcessHCLBlock(hclBlock, ctx); err != nil {
					return nil, err
				}
				dGraphNodeCount++
				instance.GraphID = dGraphNodeCount
				blockPathToGraphID[BuildBlockPath(blockName, instance.Name)] = instance.GraphID

				config.Instances[instance.Name] = instance
				dependencyGraph.AddNode(&instance)
			}
		case "data":
			config.DataBlocks["region"] = make(map[string]Block)
			config.DataBlocks["plan"] = make(map[string]Block)
			config.DataBlocks["os"] = make(map[string]Block)

			for _, hclBlock := range hclBlocks {
				dGraphNodeCount++
				dataBlock := DataBlock{
					Type:    hclBlock.Labels[0],
					Name:    hclBlock.Labels[1],
					GraphID: dGraphNodeCount,
				}
				blockPathToGraphID[BuildBlockPath(blockName, dataBlock.Type, dataBlock.Name)] = dataBlock.GraphID

				switch dataBlock.Type {
				case "region":
					regionData := RegionDataBlock{DataBlock: dataBlock}
					if err := regionData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.DataBlocks[dataBlock.Type][dataBlock.Name] = &regionData
					dependencyGraph.AddNode(&regionData)
				case "plan":
					planData := PlanDataBlock{DataBlock: dataBlock}
					if err := planData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.DataBlocks[dataBlock.Type][dataBlock.Name] = &planData
					dependencyGraph.AddNode(&planData)
				case "os":
					osData := OSDataBlock{DataBlock: dataBlock}
					if err := osData.PreProcessHCLBlock(hclBlock, ctx); err != nil {
						return nil, err
					}

					config.DataBlocks[dataBlock.Type][dataBlock.Name] = &osData
					dependencyGraph.AddNode(&osData)
				default:
					return nil, errors.New("unknown data type " + dataBlock.Type)
				}
			}
		default:
			fmt.Println("unknown block type", blockName)
		}
	}

	fmt.Println()

	fmt.Println("dependencyGraph.Nodes().Len():", dependencyGraph.Nodes().Len())

	nodes := dependencyGraph.Nodes()
	for nodes.Next() {
		node := nodes.Node().(Block)
		fmt.Println("node:", node.ID(), node)
	}

	fmt.Println()
	nodes = dependencyGraph.Nodes()
	for nodes.Next() {
		node := nodes.Node().(Block)
		fmt.Println("node:", node.ID(), node)
	}

	return &config, nil
}
