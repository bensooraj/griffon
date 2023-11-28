package parser

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/bensooraj/griffon/blocks"
	"github.com/bensooraj/griffon/mocks"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/vultr/govultr/v3"
	"go.uber.org/mock/gomock"
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
			config, err := ParseWithBodySchema("test.hcl", tC.src, evalCtx, nil)
			if diag, ok := err.(hcl.Diagnostics); ok && diag.HasErrors() {
				for i, diagErr := range diag.Errs() {
					t.Log("HCL diagnostic error [", i, "]:", diagErr.Error())
				}
			}
			require.NoError(t, err)
			require.Equalf(t, tC.expected.Griffon, config.Griffon, "GriffonBlock parsed incorrectly")
			require.Equalf(
				t,
				tC.expected.Resources[blocks.SSHKeyBlockType]["my_key"],
				config.Resources[blocks.SSHKeyBlockType]["my_key"],
				"GriffonBlock parsed incorrectly",
			)
		})
	}
}

func TestAPICall(t *testing.T) {
	t.Setenv("VULTR_API_KEY", "AxDfCdASdFzzxserDFWSD")

	defer t.Cleanup(func() {
		t.Setenv("VULTR_API_KEY", "")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVultr := mocks.NewMockVultrClient(ctrl)
	mockStartupScriptService := mocks.NewMockStartupScriptService(ctrl)
	mockStartupScriptService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&govultr.StartupScript{
		ID:           "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		DateCreated:  "2020-10-10T01:56:20+00:00",
		DateModified: "2020-10-10T01:59:20+00:00",
		Name:         "my_key",
		Type:         "pxe",
		Script:       "ssh-rsa AAAAB3NzaC1yc2E",
	}, &http.Response{}, nil)

	mockVultr.StartupScript = mockStartupScriptService

	b, err := os.ReadFile("../testdata/test_1.hcl")
	if err != nil {
		panic(err)
	}
	// parse the file
	config, err := ParseWithBodySchema("testdata/test1.hcl", b, GetEvalContext(), mockVultr)
	require.NoError(t, err)

	// check if the parsed config is correct
	require.Equalf(
		t,
		blocks.GriffonBlock{Region: "us-east-1", VultrAPIKey: "AxDfCdASdFzzxserDFWSD"},
		config.Griffon,
		"GriffonBlock parsed incorrectly",
	)
}
