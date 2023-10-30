package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGriffonBlock(t *testing.T) {
	src := []byte(`
	griffon {
		region = "nyc1"
		vultr_api_key = "1234567890"
	}`)

	config, err := ParseHCL("test.hcl", src)
	require.NoError(t, err)
	require.Equalf(t, GriffonBlock{
		Region:      "nyc1",
		VultrAPIKey: "1234567890",
	}, config.Griffon, "GriffonBlock parsed incorrectly")
}

func TestParseSshKeyBlock(t *testing.T) {
	src := []byte(`
	griffon {
		region = "nyc1"
		vultr_api_key = "1234567890"
	}
	
	ssh_key "my_key" {
		ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADA"
	}`)

	config, err := ParseHCL("test.hcl", src)
	require.NoError(t, err)
	require.Equalf(t, GriffonBlock{
		Region:      "nyc1",
		VultrAPIKey: "1234567890",
	}, config.Griffon, "GriffonBlock parsed incorrectly")

	require.NotNil(t, config.SSHKeys, "SSHKeys not parsed")

	require.Lenf(t, config.SSHKeys, 1, "len(SSHKeys); got %d, want 1", len(config.SSHKeys))

	require.Equalf(t, SSHKeyBlock{
		Name:   "my_key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA",
	}, config.SSHKeys[0], "SSHKeyBlock parsed incorrectly")
}
