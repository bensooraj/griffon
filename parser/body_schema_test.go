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

	mockVultr := setupMockVultrClient(t, ctrl)

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

	t.Logf("config.Resources: %+v", config.Resources)
	expectedSSHKeyBlock := blocks.SSHKeyBlock{SSHKey: "ssh-rsa AAAAB3NzaC1yc2E", ResourceBlock: blocks.ResourceBlock{Name: "my_key", Type: "ssh_key"}}
	require.Equalf(t, expectedSSHKeyBlock, config.Resources[blocks.SSHKeyBlockType]["my_key"], "SSHKeyBlock parsed incorrectly")

}

func TestEvaluateConfig(t *testing.T) {

	t.Setenv("VULTR_API_KEY", "AxDfCdASdFzzxserDFWSD")

	defer t.Cleanup(func() {
		t.Setenv("VULTR_API_KEY", "")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVultr := setupMockVultrClient(t, ctrl)

	type args struct {
		config *blocks.Config
		vc     *govultr.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test evaluate config",
			args:    args{config: &blocks.Config{}, vc: mockVultr},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EvaluateConfig(tt.args.config, tt.args.vc); (err != nil) != tt.wantErr {
				t.Errorf("EvaluateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setupMockVultrClient(t *testing.T, ctrl *gomock.Controller) *govultr.Client {
	t.Helper()

	mockVultr := mocks.NewMockVultrClient(ctrl)

	// 1. Data
	// 1.1 Regions
	mockRegionService := mocks.NewMockRegionService(ctrl)
	mockRegionService.EXPECT().List(gomock.Any(), gomock.Any()).Return([]govultr.Region{
		{
			ID:        "ams",
			City:      "Amsterdam",
			Country:   "NL",
			Continent: "Europe",
			Options:   []string{"ddos_protection", "block_storage_high_perf", "block_storage_storage_opt", "kubernetes", "load_balancers"},
		},
	}, &govultr.Meta{}, &http.Response{}, nil).AnyTimes()
	mockVultr.Region = mockRegionService
	// 1.2 OS
	mockOSService := mocks.NewMockOSService(ctrl)
	mockOSService.EXPECT().List(gomock.Any(), gomock.Any()).Return([]govultr.OS{}, &govultr.Meta{}, &http.Response{}, nil).AnyTimes()
	mockVultr.OS = mockOSService
	// 1.3 Plans
	mockPlanService := mocks.NewMockPlanService(ctrl)
	mockPlanService.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]govultr.Plan{
		{
			ID:          "vhf-8c-32gb",
			VCPUCount:   8,
			RAM:         32768,
			Disk:        512,
			DiskCount:   1,
			Bandwidth:   6144,
			MonthlyCost: 192,
			Type:        "vhf",
			Locations:   []string{"ams", "atl", "dfw", "fra", "hnd", "lax", "lhr", "mia", "nrt", "ord", "sea", "sgp", "sjc", "syd", "tor"},
		},
	}, &govultr.Meta{}, &http.Response{}, nil).AnyTimes()
	mockVultr.Plan = mockPlanService

	// 2. Resources
	// 2.1 Instance
	mockStartupScriptService := mocks.NewMockStartupScriptService(ctrl)
	mockStartupScriptService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&govultr.StartupScript{
		ID:           "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		DateCreated:  "2020-10-10T01:56:20+00:00",
		DateModified: "2020-10-10T01:59:20+00:00",
		Name:         "my_key",
		Type:         "pxe",
		Script:       "ssh-rsa AAAAB3NzaC1yc2E",
	}, &http.Response{}, nil).AnyTimes()
	mockVultr.StartupScript = mockStartupScriptService
	// 2.2 SSHKey
	mockSSHKeyService := mocks.NewMockSSHKeyService(ctrl)
	mockSSHKeyService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&govultr.SSHKey{
		ID:          "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		Name:        "my_key",
		SSHKey:      "ssh-rsa AAAAB3NzaC1yc2E",
		DateCreated: "2020-10-10T01:56:20+00:00",
	}, &http.Response{}, nil).AnyTimes()
	mockVultr.SSHKey = mockSSHKeyService
	// 2.3 Instance
	mockInstanceService := mocks.NewMockInstanceService(ctrl)
	mockInstanceService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&govultr.Instance{
		ID:               "4f0f12e5-1f84-404f-aa84-85f431ea5ec2",
		Region:           "ams",
		Plan:             "vc2-1c-1gb",
		Os:               "CentOS 8 Stream",
		OsID:             362,
		RAM:              1024,
		Disk:             25,
		MainIP:           "",
		VCPUCount:        1,
		DateCreated:      "2020-10-10T01:56:20+00:00",
		Status:           "pending",
		AllowedBandwidth: 1000,
		DefaultPassword:  "v5{Fkvb#2ycPGwHs",
		NetmaskV4:        "",
		GatewayV4:        "0.0.0.0",
		PowerStatus:      "running",
		ServerStatus:     "ok",
		Tags: []string{
			"tag1",
		},
	}, &http.Response{}, nil).AnyTimes()
	mockVultr.Instance = mockInstanceService

	return mockVultr
}

func TestEvaluationContext(t *testing.T) {
	evalCtx := GetEvalContext()
	//
	err := AddBlockToEvalContext(evalCtx, &blocks.RegionDataBlock{
		DataBlock: blocks.DataBlock{Type: blocks.RegionBlockType, Name: "region_1"},
		VID:       "ams",
		City:      "Amsterdam",
		Country:   "NL",
		Continent: "Europe",
	})
	if err != nil {
		t.Log("Error adding region data block to eval context:", err)
	}

	// add another region
	err = AddBlockToEvalContext(evalCtx, &blocks.RegionDataBlock{
		DataBlock: blocks.DataBlock{Type: blocks.RegionBlockType, Name: "region_2"},
		VID:       "use",
		City:      "New York",
		Country:   "US",
		Continent: "North America",
	})
	if err != nil {
		t.Log("Error adding another region data block to eval context:", err)
	}

	// add a plan
	err = AddBlockToEvalContext(evalCtx, &blocks.PlanDataBlock{
		DataBlock:   blocks.DataBlock{Type: blocks.PlanBlockType, Name: "plan_1"},
		VID:         "vhf-8c-32gb",
		VCPUCount:   8,
		RAM:         32768,
		Disk:        512,
		DiskCount:   1,
		Bandwidth:   6144,
		MonthlyCost: 192,
		PlanType:    "vhf",
	})
	if err != nil {
		t.Log("Error adding plan data block to eval context:", err)
	}

	// add an OS
	err = AddBlockToEvalContext(evalCtx, &blocks.OSDataBlock{
		DataBlock: blocks.DataBlock{Type: blocks.OSBlockType, Name: "os_1"},
		VID:       362,
		OSName:    "CentOS 8 Stream",
		Arch:      "x64",
		Family:    "centos",
	})
	if err != nil {
		t.Log("Error adding OS data block to eval context:", err)
	}

	// add an SSHKey
	err = AddBlockToEvalContext(evalCtx, &blocks.SSHKeyBlock{
		ResourceBlock: blocks.ResourceBlock{Type: blocks.SSHKeyBlockType, Name: "ssh_key_1"},
		VID:           "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		SSHKey:        "ssh-rsa AAAAB3NzaC1yc2E",
		DateCreated:   "2020-10-10T01:56:20+00:00",
	})
	if err != nil {
		t.Log("Error adding SSH Key resource block to eval context:", err)
	}

	// add another SSHKey
	err = AddBlockToEvalContext(evalCtx, &blocks.SSHKeyBlock{
		ResourceBlock: blocks.ResourceBlock{Type: blocks.SSHKeyBlockType, Name: "ssh_key_2"},
		VID:           "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		SSHKey:        "ssh-rsa AAAAB3NzaC1yc2E",
		DateCreated:   "2020-10-10T01:56:20+00:00",
	})
	if err != nil {
		t.Log("Error adding another SSH Key resource block to eval context:", err)
	}

	// add a StartupScript
	err = AddBlockToEvalContext(evalCtx, &blocks.StartupScriptBlock{
		ResourceBlock: blocks.ResourceBlock{Type: blocks.StartupScriptBlockType, Name: "startup_script_1"},
		VID:           "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		VDateCreated:  "2020-10-10T01:56:20+00:00",
		VDateModified: "2020-10-10T01:59:20+00:00",
		VType:         "pxe",
	})
	if err != nil {
		t.Log("Error adding StartupScript resource block to eval context:", err)
	}

	// add an Instance
	err = AddBlockToEvalContext(evalCtx, &blocks.InstanceBlock{
		ResourceBlock: blocks.ResourceBlock{Type: blocks.InstanceBlockType, Name: "instance_1"},
		VID:           "4f0f12e5-1f84-404f-aa84-85f431ea5ec2",
		Os:            "CentOS 8 Stream",
		RAM:           1024,
		Disk:          25,
		VCPUCount:     1,
		Status:        "pending",
	})
	if err != nil {
		t.Log("Error adding Instance resource block to eval context:", err)
	}

	fmt.Println("[AFTER] data:", evalCtx.Variables["data"].GoString())
	fmt.Println("[AFTER] SSH Key:", evalCtx.Variables[string(blocks.SSHKeyBlockType)].GoString())
	fmt.Println("[AFTER] Startup Script:", evalCtx.Variables[string(blocks.StartupScriptBlockType)].GoString())
	fmt.Println("[AFTER] Instance:", evalCtx.Variables[string(blocks.InstanceBlockType)].GoString())

	t.Fail()
}
