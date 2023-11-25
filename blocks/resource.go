package blocks

import "github.com/hashicorp/hcl/v2"

type ResourceBlock struct {
	GraphID   int64  `json:"graph_id"`
	Type      string `hcl:"type,label"`
	Name      string `hcl:"name,label"`
	Config    hcl.Body
	DependsOn []string `json:"depends_on"`
}
