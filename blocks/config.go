package blocks

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
	Griffon   GriffonBlock                   `hcl:"griffon,block" json:"griffon"`
	Data      map[BlockType]map[string]Block `hcl:"data,block" json:"data"`
	Resources map[BlockType]map[string]Block `hcl:"resources,block" json:"resources"`
}

func (c *Config) AddResource(b Block) {
	if _, ok := c.Resources[b.BlockType()]; !ok {
		c.Resources[b.BlockType()] = make(map[string]Block)
	}
	c.Resources[b.BlockType()][b.BlockName()] = b
}
