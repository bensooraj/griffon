package blocks

import "github.com/hashicorp/hcl/v2"

type ResourceBlock struct {
	GraphID   int64
	Type      string `hcl:"type,label"`
	Name      string `hcl:"name,label"`
	Config    hcl.Body
	DependsOn []string
}
