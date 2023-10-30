package main

type Config struct {
	Griffon GriffonBlock  `hcl:"griffon,block"`
	SSHKeys []SSHKeyBlock `hcl:"ssh_key,block"`
}

type GriffonBlock struct {
	Region      string `hcl:"region,attr"`
	VultrAPIKey string `hcl:"vultr_api_key"`
}

type SSHKeyBlock struct {
	Name   string `hcl:"name,label"`
	SSHKey string `hcl:"ssh_key"`
}
