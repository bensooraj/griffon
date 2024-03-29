package blocks

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type SSHKeyBlock struct {
	SSHKey      string `hcl:"ssh_key" json:"ssh_key,omitempty" cty:"ssh_key"`
	DateCreated string `json:"date_created,omitempty" cty:"date_created"`
	VID         string `json:"id,omitempty" cty:"id"`
	ResourceBlock
}

var _ Block = (*SSHKeyBlock)(nil)

func (s *SSHKeyBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := s.Config.PartialContent(schema.SSHKeyBlockSchema)
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
	slog.Debug("SSH Key parameters", slog.String("block_type", string(s.BlockType())), slog.String("block_name", string(s.BlockName())), slog.Any("params", s))

	return nil
}

func (s *SSHKeyBlock) Create(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) error {
	sshKey, _, err := vc.SSHKey.Create(context.Background(), &govultr.SSHKeyReq{
		Name:   s.Name,
		SSHKey: s.SSHKey,
	})
	if err != nil {
		return err
	}

	s.VID = sshKey.ID
	s.DateCreated = sshKey.DateCreated
	slog.Info("SSH Key created", slog.String("block_type", string(s.BlockType())), slog.String("block_name", string(s.BlockName())), slog.Any("ssh_key", sshKey))

	return nil
}

// ToCtyValue
func (s *SSHKeyBlock) ToCtyValue() (cty.Value, error) {
	return gocty.ToCtyValue(s, cty.Object(map[string]cty.Type{
		"id":           cty.String,
		"ssh_key":      cty.String,
		"date_created": cty.String,
	}))
}
