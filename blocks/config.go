package blocks

import "log/slog"

type BlockType string

const (
	GriffonBlockType BlockType = "griffon"

	DataBlockType   BlockType = "data"
	RegionBlockType BlockType = "region"
	PlanBlockType   BlockType = "plan"
	OSBlockType     BlockType = "os"

	ResourceBlockType      BlockType = "resource"
	SSHKeyBlockType        BlockType = "ssh_key"
	StartupScriptBlockType BlockType = "startup_script"
	InstanceBlockType      BlockType = "instance"
)

type Config struct {
	EvaluationOrder []int64                        `json:"evaluation_order,omitempty"`
	Griffon         GriffonBlock                   `hcl:"griffon,block" json:"griffon"`
	Data            map[BlockType]map[string]Block `hcl:"data,block" json:"data"`
	Resources       map[BlockType]map[string]Block `hcl:"resources,block" json:"resources"`
}

func (c *Config) AddResource(b Block) {
	if _, ok := c.Resources[b.BlockType()]; !ok {
		c.Resources[b.BlockType()] = make(map[string]Block)
	}
	c.Resources[b.BlockType()][b.BlockName()] = b
	slog.Debug("Added resource to config", slog.String("block_type", string(b.BlockType())), slog.String("block_name", b.BlockName()))
}

func (c *Config) AddData(b Block) {
	if _, ok := c.Data[b.BlockType()]; !ok {
		c.Data[b.BlockType()] = make(map[string]Block)
	}
	c.Data[b.BlockType()][b.BlockName()] = b
	slog.Debug("Added data to config", slog.String("block_type", string(b.BlockType())), slog.String("block_name", b.BlockName()))
}
