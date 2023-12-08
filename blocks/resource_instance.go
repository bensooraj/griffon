package blocks

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/bensooraj/griffon/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/vultr/govultr/v3"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type InstanceBlock struct {
	Region          string   `hcl:"region,attr" json:"region,omitempty"`
	Plan            string   `hcl:"plan,attr" json:"plan,omitempty"`
	OsID            int      `hcl:"os_id,attr" json:"os_id,omitempty"`
	SshKeyID        string   `hcl:"ssh_key_id,attr" json:"sshkey_id,omitempty"`
	StartupScriptID string   `hcl:"script_id,attr" json:"script_id,omitempty"`
	Hostname        string   `hcl:"hostname,attr" json:"hostname"`
	Tags            []string `hcl:"tag,attr" json:"tags"`

	VID             string `json:"id" cty:"id"`
	Os              string `json:"os" cty:"os"`
	RAM             int    `json:"ram" cty:"ram"`
	Disk            int    `json:"disk" cty:"disk"`
	MainIP          string `json:"main_ip" cty:"main_ip"`
	VCPUCount       int    `json:"vcpu_count" cty:"vcpu_count"`
	DefaultPassword string `json:"default_password,omitempty" cty:"default_password"`
	DateCreated     string `json:"date_created" cty:"date_created"`
	Status          string `json:"status" cty:"status"`
	PowerStatus     string `json:"power_status" cty:"power_status"`
	ServerStatus    string `json:"server_status" cty:"server_status"`
	InternalIP      string `json:"internal_ip" cty:"internal_ip"`

	ResourceBlock
}

var _ Block = (*InstanceBlock)(nil)

func (i *InstanceBlock) ProcessConfiguration(ctx *hcl.EvalContext) error {
	content, _, diags := i.Config.PartialContent(schema.InstanceBlockSchema)
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
			osID, err := strconv.Atoi(value.AsString())
			if err != nil {
				return err
			}
			i.OsID = osID
		case "ssh_key_id":
			i.SshKeyID = value.AsString()
		case "startup_script_id":

		case "hostname":
			i.Hostname = value.AsString()
		case "tag":
			fmt.Println("tags:", value.AsString(), value.AsValueMap())
			for key, ctyVal := range value.AsValueMap() {
				i.Tags = append(i.Tags, fmt.Sprintf("%s=%s", key, ctyVal.AsString()))
			}
		}
	}

	return nil
}

func (i *InstanceBlock) Create(ctx context.Context, evalCtx *hcl.EvalContext, vc *govultr.Client) (*hcl.EvalContext, error) {
	fmt.Println("Creating instance", i.Name)
	ins, _, err := vc.Instance.Create(ctx, &govultr.InstanceCreateReq{
		Region:   i.Region,
		Plan:     i.Plan,
		OsID:     i.OsID,
		SSHKeys:  []string{i.SshKeyID},
		ScriptID: i.StartupScriptID,
		Hostname: i.Hostname,
		Tags:     i.Tags,
	})
	if err != nil {
		return nil, err
	}

	i.VID = ins.ID
	i.Os = ins.Os
	i.RAM = ins.RAM
	i.Disk = ins.Disk
	i.MainIP = ins.MainIP
	i.VCPUCount = ins.VCPUCount
	i.DefaultPassword = ins.DefaultPassword
	i.DateCreated = ins.DateCreated
	i.Status = ins.Status
	i.PowerStatus = ins.PowerStatus
	i.ServerStatus = ins.ServerStatus
	i.InternalIP = ins.InternalIP

	fmt.Println("Created instance", i.Name, i.VID)
	fmt.Printf("Created instance %+v\n", ins)

	return nil, nil
}

func (i *InstanceBlock) ToCtyValue() (cty.Value, error) {
	return gocty.ToCtyValue(i, cty.Object(map[string]cty.Type{
		"id":               cty.String,
		"os":               cty.String,
		"ram":              cty.Number,
		"disk":             cty.Number,
		"main_ip":          cty.String,
		"vcpu_count":       cty.Number,
		"default_password": cty.String,
		"date_created":     cty.String,
		"status":           cty.String,
		"power_status":     cty.String,
		"server_status":    cty.String,
		"internal_ip":      cty.String,
	}))
}
