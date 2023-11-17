package main

import "github.com/hashicorp/hcl/v2"

type Config struct {
	Griffon        GriffonBlock         `hcl:"griffon,block"`
	SSHKeys        []SSHKeyBlock        `hcl:"ssh_key,block"`
	StartupScripts []StartupScriptBlock `hcl:"startup_script,block"`
}

type GriffonBlock struct {
	Region      string `hcl:"region,attr"`
	VultrAPIKey string `hcl:"vultr_api_key"`
}

type SSHKeyBlock struct {
	Name   string `hcl:"name,label"`
	SSHKey string `hcl:"ssh_key"`
}

type StartupScriptBlock struct {
	Name      string `hcl:"name,label"`
	Script    string `hcl:"script"`
	Config    hcl.Body
	DependsOn map[string][]string
}

type DataBlock struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}

type InstanceBlock struct {
	Name            string            `hcl:"name,label"`
	Region          string            `hcl:"region,attr"`
	Plan            string            `hcl:"plan,attr"`
	OS              string            `hcl:"os,attr"`
	SshKeyID        string            `hcl:"ssh_key_id,attr"`
	StartupScriptID string            `hcl:"startup_script_id,attr"`
	Hostname        string            `hcl:"hostname,attr"`
	Tag             map[string]string `hcl:"tag,attr"`
}
