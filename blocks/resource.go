package blocks

import "github.com/hashicorp/hcl/v2"

type ResourceBlock struct {
	GraphID   int64  `json:"graph_id"`
	Type      string `hcl:"type,label" json:"type"`
	Name      string `hcl:"name,label" json:"name"`
	Config    hcl.Body
	DependsOn []string `json:"depends_on"`
}

func (r *ResourceBlock) ID() int64 {
	return r.GraphID
}

func (r *ResourceBlock) Dependencies() []string {
	return r.DependsOn
}
