package main

import "github.com/hashicorp/hcl/v2"

type Config struct {
	Griffon        GriffonBlock                  `hcl:"griffon,block"`
	SSHKeys        map[string]SSHKeyBlock        `hcl:"ssh_key,block"`
	StartupScripts map[string]StartupScriptBlock `hcl:"startup_script,block"`
	Instances      map[string]InstanceBlock      `hcl:"instance,block"`
	DataBlocks     map[string]map[string]Block   `hcl:"data,block"`
}

type GriffonBlock struct {
	GraphID     int64
	Region      string `hcl:"region,attr"`
	VultrAPIKey string `hcl:"vultr_api_key"`
}

type SSHKeyBlock struct {
	GraphID   int64
	Name      string `hcl:"name,label"`
	SSHKey    string `hcl:"ssh_key"`
	Config    hcl.Body
	DependsOn map[string][]string
}

type StartupScriptBlock struct {
	GraphID   int64
	Name      string `hcl:"name,label"`
	Script    string `hcl:"script"`
	Config    hcl.Body
	DependsOn map[string][]string
}

type DataBlock struct {
	GraphID   int64
	Type      string `hcl:"type,label"`
	Name      string `hcl:"name,label"`
	Config    hcl.Body
	DependsOn map[string][]string
}

type RegionDataBlock struct {
	VultrID   string   `json:"id"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Continent string   `json:"continent"`
	Options   []string `json:"options"`
	DataBlock
}

type PlanDataBlock struct {
	VultrID     string   `json:"id"`
	VcpuCount   int      `json:"vcpu_count"`
	Ram         int      `json:"ram"`
	Disk        int      `json:"disk"`
	DiskCount   int      `json:"disk_count"`
	Bandwidth   int      `json:"bandwidth"`
	MonthlyCost int      `json:"monthly_cost"`
	PlanType    string   `json:"type"`
	Locations   []string `json:"locations"`
	DataBlock
}

type OSDataBlock struct {
	VultrID int    `json:"id"`
	OSName  string `json:"name"`
	Arch    string `json:"arch"`
	Family  string `json:"family"`
	DataBlock
}

type InstanceBlock struct {
	GraphID         int64
	Name            string            `hcl:"name,label"`
	Region          string            `hcl:"region,attr"`
	Plan            string            `hcl:"plan,attr"`
	OS              string            `hcl:"os,attr"`
	SshKeyID        string            `hcl:"ssh_key_id,attr"`
	StartupScriptID string            `hcl:"startup_script_id,attr"`
	Hostname        string            `hcl:"hostname,attr"`
	Tag             map[string]string `hcl:"tag,attr"`
	Config          hcl.Body
	DependsOn       map[string][]string
}

type Block interface {
	ID() int64
	// Separate the block into its configuration and dependencies
	PreProcessHCLBlock(block *hcl.Block, ctx *hcl.EvalContext) error
	// Process the configuration
	ProcessConfiguration(ctx *hcl.EvalContext) error
	// Execute the block by making API calls
	// Execute(ctx *hcl.EvalContext) error
}
