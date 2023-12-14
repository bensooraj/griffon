package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/bensooraj/griffon/blocks"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
)

//nolint:dupl
func TestSpecParser(t *testing.T) {
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
		expected blocks.Config
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
			expected: blocks.Config{
				Griffon: blocks.GriffonBlock{Region: "nyc1", VultrAPIKey: "1234567890"},
				Resources: map[blocks.BlockType]map[string]blocks.Block{
					blocks.SSHKeyBlockType: {
						"my_key": &blocks.SSHKeyBlock{
							ResourceBlock: blocks.ResourceBlock{
								Name: "my_key",
								Type: "ssh_key",
							},
							SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA",
						},
					},
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
			expected: blocks.Config{
				Griffon: blocks.GriffonBlock{Region: "AMS", VultrAPIKey: "axdfcdasdfzzxserdfwsd"},
				Resources: map[blocks.BlockType]map[string]blocks.Block{
					blocks.SSHKeyBlockType: {
						"my_key": &blocks.SSHKeyBlock{
							ResourceBlock: blocks.ResourceBlock{
								Name: "my_key",
								Type: "ssh_key",
							},
							SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA",
						},
					},
				},
			},
		},
	}

	evalCtx := GetEvalContext()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config, err := ParseHCLUsingSpec("test.hcl", tC.src, evalCtx)
			if diag, ok := err.(hcl.Diagnostics); ok && diag.HasErrors() {
				for i, diagErr := range diag.Errs() {
					t.Log("HCL diagnostic error [", i, "]:", diagErr.Error())
				}
			}
			require.NoError(t, err)
			require.Equalf(t, tC.expected.Griffon, config.Griffon, "GriffonBlock parsed incorrectly")
			require.Equalf(
				t,
				tC.expected.Resources[blocks.SSHKeyBlockType]["my_key"].(*blocks.SSHKeyBlock).SSHKey,
				config.Resources[blocks.SSHKeyBlockType]["my_key"].(*blocks.SSHKeyBlock).SSHKey,
				"GriffonBlock parsed incorrectly",
			)
		})
	}
}
