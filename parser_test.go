package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
)

func TestParseGriffonBlock(t *testing.T) {
	src := []byte(`
	griffon {
		region = "nyc1"
		vultr_api_key = "1234567890"
	}`)

	config, err := ParseHCL("test.hcl", src, nil)
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

	config, err := ParseHCL("test.hcl", src, nil)
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

func TestParseGriffonBlock_Variables(t *testing.T) {

	t.Setenv("VULTR_API_KEY", "1234567890")
	defer t.Cleanup(func() {
		t.Setenv("VULTR_API_KEY", "")
	})

	testCases := []struct {
		desc     string
		src      []byte
		expected GriffonBlock
	}{
		{
			desc: "parse AMS as a variable",
			src: []byte(`
			griffon {
				region = AMS
				vultr_api_key = "1234567890"
			}`),
			expected: GriffonBlock{Region: "ams", VultrAPIKey: "1234567890"},
		},
		{
			desc: "parse AMS as a template string",
			src: []byte(`
			griffon {
				region = "${AMS}terdam"
				vultr_api_key = "1234567890"
			}`),
			expected: GriffonBlock{Region: "amsterdam", VultrAPIKey: "1234567890"},
		},
		{
			desc: "parse AMS as a template string",
			src: []byte(`
			griffon {
				region = "${AMS == "ams" ? "toronto" : "amsterdam"}"
				vultr_api_key = "1234567890"
			}`),
			expected: GriffonBlock{Region: "toronto", VultrAPIKey: "1234567890"},
		},
		{
			desc: "parse AMS as a variable",
			src: []byte(`
			griffon {
				region = AMS
				vultr_api_key = env.VULTR_API_KEY
			}`),
			expected: GriffonBlock{Region: "ams", VultrAPIKey: "1234567890"},
		},
	}
	evalCtx := getEvalContext()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config, err := ParseHCL("test.hcl", tC.src, evalCtx)
			if diag, ok := err.(hcl.Diagnostics); ok && diag.HasErrors() {
				for i, diagErr := range diag.Errs() {
					t.Log("HCL diagnostic error [", i, "]:", diagErr.Error())
				}
			}
			require.NoError(t, err)
			require.Equalf(t, tC.expected, config.Griffon, "GriffonBlock parsed incorrectly")
		})
	}
}

func Test5_Functions(t *testing.T) {

	t.Setenv("VULTR_API_KEY", "AxDfCdASdFzzxserDFWSD")
	myKeyPubFile, err := os.CreateTemp("", "my_key.pub")
	require.NoError(t, err)
	_, err = myKeyPubFile.WriteString("ssh-rsa AAAAB3NzaC1yc2EAAAADA")
	require.NoError(t, err)

	defer t.Cleanup(func() {
		t.Setenv("VULTR_API_KEY", "")
		os.Remove(myKeyPubFile.Name())
	})

	testCases := []struct {
		desc     string
		src      []byte
		expected Config
	}{
		{
			desc: "built-in functions uppercase and lowercase",
			src: []byte(`
			griffon {
				region = uppercase(AMS)
				vultr_api_key = lowercase(env.VULTR_API_KEY)
			}
			ssh_key "my_key" {
				ssh_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADA"
			}`),
			expected: Config{
				Griffon: GriffonBlock{Region: "AMS", VultrAPIKey: "axdfcdasdfzzxserdfwsd"},
				SSHKeys: []SSHKeyBlock{
					{Name: "my_key", SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA"},
				},
			},
		},
		{
			desc: "custom functions file",
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
				SSHKeys: []SSHKeyBlock{
					{Name: "my_key", SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADA"},
				},
			},
		},
	}
	evalCtx := getEvalContext()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config, err := ParseHCL("test.hcl", tC.src, evalCtx)
			if diag, ok := err.(hcl.Diagnostics); ok && diag.HasErrors() {
				for i, diagErr := range diag.Errs() {
					t.Log("HCL diagnostic error [", i, "]:", diagErr.Error())
				}
			}
			require.NoError(t, err)
			require.Equalf(t, tC.expected.Griffon, config.Griffon, "GriffonBlock parsed incorrectly")
			require.Equalf(t, tC.expected.SSHKeys[0], config.SSHKeys[0], "GriffonBlock parsed incorrectly")
		})
	}
}
