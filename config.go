package main

type Config struct {
	Griffon GriffonBlock `hcl:"griffon,block"`
}

type GriffonBlock struct {
	Region      string `hcl:"region,attr"`
	VultrAPIKey string `hcl:"vultr_api_key"`
}
