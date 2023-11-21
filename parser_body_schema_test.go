package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
)

//nolint:dupl
func TestBodySchemaParser(t *testing.T) {
	t.Setenv("VULTR_API_KEY", "AxDfCdASdFzzxserDFWSD")
	myKeyPubFile, err := os.CreateTemp("", "my_key.pub")
	require.NoError(t, err)
	_, err = myKeyPubFile.WriteString("ssh-rsa AAAAB3NzaC1yc2EAAAADA")
	require.NoError(t, err)

	defer t.Cleanup(func() {
		t.Setenv("VULTR_API_KEY", "")
		err := os.Remove(myKeyPubFile.Name())
		if err != nil {
			t.Log("Error removing temp file:", err)
		}
	})

	testCases := []struct {
		desc     string
		src      []byte
		expected Config
	}{
		{
			desc: "simple config",
			src: []byte(`
			griffon {
				region = "nyc1"
				vultr_api_key = "1234567890"
			}
			ssh_key "my_key" {
				ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADA"
			}`),
			expected: Config{
				Griffon: GriffonBlock{Region: "nyc1", VultrAPIKey: "1234567890"},
				SSHKeys: map[string]SSHKeyBlock{
					"my_key": {Name: "my_key", SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA"},
				},
			},
		},
		{
			desc: "parse HCL with variables and functions",
			src: []byte(fmt.Sprintf(`
			griffon {
				region = uppercase(AMS)
				vultr_api_key = lowercase(env.VULTR_API_KEY)
			}
			ssh_key "my_key" {
				ssh_key = file("%s")
			}`, myKeyPubFile.Name())),
			expected: Config{
				Griffon: GriffonBlock{Region: "AMS", VultrAPIKey: "axdfcdasdfzzxserdfwsd"},
				SSHKeys: map[string]SSHKeyBlock{
					"my_key": {Name: "my_key", SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA"},
				},
			},
		},
	}

	evalCtx := getEvalContext()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config, err := ParseHCLUsingBodySchema("test.hcl", tC.src, evalCtx)
			if diag, ok := err.(hcl.Diagnostics); ok && diag.HasErrors() {
				for i, diagErr := range diag.Errs() {
					t.Log("HCL diagnostic error [", i, "]:", diagErr.Error())
				}
			}
			require.NoError(t, err)
			require.Equalf(t, tC.expected.Griffon, config.Griffon, "GriffonBlock parsed incorrectly")
			require.Equalf(t, tC.expected.SSHKeys["my_key"], config.SSHKeys["my_key"], "GriffonBlock parsed incorrectly")
		})
	}
}
