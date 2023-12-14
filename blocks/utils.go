package blocks

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func ExprAsMap(expr hcl.Expression) (map[string][]string, hcl.Diagnostics) {
	exprMap := make(map[string][]string)

	exprList, diags := hcl.ExprList(expr)
	if diags.HasErrors() {
		return nil, diags
	}
	for _, v := range exprList {
		traversals := v.Variables()
		for _, traversal := range traversals {
			for _, step := range traversal {
				switch t := step.(type) {
				case hcl.TraverseAttr:
					exprMap[traversal.RootName()] = append(exprMap[traversal.RootName()], t.Name)
				case hcl.TraverseIndex:
				case hcl.TraverseRoot:
				default:
				}
			}
		}
	}
	return exprMap, nil
}

func ExprAsStringSlice(expr hcl.Expression) ([]string, hcl.Diagnostics) {
	var exprSlice []string

	exprList, diags := hcl.ExprList(expr)
	if diags.HasErrors() {
		return nil, diags
	}
	for _, v := range exprList {
		traversals := v.Variables()
		for _, traversal := range traversals {
			var paths []string
			for _, step := range traversal {
				switch t := step.(type) {
				case hcl.TraverseRoot:
					paths = append(paths, t.Name)
				case hcl.TraverseAttr:
					paths = append(paths, t.Name)
				case hcl.TraverseIndex:
				default:
				}
			}
			exprSlice = append(exprSlice, BuildBlockPath(paths...))
		}
	}
	return exprSlice, nil
}

func BuildBlockPath(paths ...string) string {
	return strings.Join(paths, ".")
}
