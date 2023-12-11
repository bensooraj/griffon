package blocks

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type StartupScriptBlock struct {
	Script        string `hcl:"script" json:"script" cty:"script"`
	VID           string `json:"id" cty:"id"`
	VDateCreated  string `json:"date_created" cty:"date_created"`
	VDateModified string `json:"date_modified" cty:"date_modified"`
	VType         string `json:"type" cty:"type"`
	ResourceBlock
}

var _ Block = (*StartupScriptBlock)(nil)

func (s *StartupScriptBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := s.Config.PartialContent(schema.StartupScriptBlockSchema)
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
			fmt.Println("script:", value.AsString())
			s.Script = base64.StdEncoding.EncodeToString([]byte(value.AsString()))
		default:
			return errors.New("unknown attribute " + attrName)
		}
	}
	return nil
}

func (s *StartupScriptBlock) Create(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error {
	fmt.Println("Creating startup script", s.Name)
	ss, _, err := vc.StartupScript.Create(context.Background(), &govultr.StartupScriptReq{
		Name:   s.Name,
		Script: s.Script,
	})
	if err != nil {
		return err
	}

	s.VID = ss.ID
	s.VDateCreated = ss.DateCreated
	s.VDateModified = ss.DateModified
	s.VType = ss.Type

	return nil
}

// ToCtyValue
func (s *StartupScriptBlock) ToCtyValue() (cty.Value, error) {
	return gocty.ToCtyValue(s, cty.Object(map[string]cty.Type{
		"id":            cty.String,
		"script":        cty.String,
		"date_created":  cty.String,
		"date_modified": cty.String,
		"type":          cty.String,
	}))
}
