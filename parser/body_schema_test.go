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
			config, err := ParseWithBodySchema("test.hcl", tC.src, evalCtx)
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

func TestParseWithBodySchema(t *testing.T) {
	t.Setenv("VULTR_API_KEY", "AxDfCdASdFzzxserDFWSD")

	defer t.Cleanup(func() {
		t.Setenv("VULTR_API_KEY", "")
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b, err := os.ReadFile("../testdata/test_1.hcl")
	if err != nil {
		panic(err)
	}
	// parse the file
	config, err := ParseWithBodySchema("testdata/test1.hcl", b, GetEvalContext())
	require.NoError(t, err)

	// GriffonBlock
	require.Equalf(
		t,
		blocks.GriffonBlock{Region: "ams", VultrAPIKey: "AxDfCdASdFzzxserDFWSD"},
		config.Griffon,
		"GriffonBlock parsed incorrectly",
	)

	t.Logf("config.Resources: %+v", config.Resources)

	// RegionDataBlock
	expectedRegionDataBlock := blocks.RegionDataBlock{VID: "ams", City: "Amsterdam", Country: "NL", Continent: "Europe", DataBlock: blocks.DataBlock{Type: blocks.RegionBlockType, Name: "current"}}
	actualRegionDataBlock := config.Data[blocks.RegionBlockType]["current"].(*blocks.RegionDataBlock)
	require.Equalf(t, blocks.RegionBlockType, actualRegionDataBlock.BlockType(), "RegionDataBlock BlockType() incorrect")
	require.Equalf(t, expectedRegionDataBlock.Name, actualRegionDataBlock.BlockName(), "RegionDataBlock BlockName() incorrect")
	require.Equalf(t, expectedRegionDataBlock.VID, actualRegionDataBlock.VID, "RegionDataBlock VID incorrect")
	require.Equalf(t, expectedRegionDataBlock.City, actualRegionDataBlock.City, "RegionDataBlock City incorrect")
	require.Equalf(t, expectedRegionDataBlock.Country, actualRegionDataBlock.Country, "RegionDataBlock Country incorrect")
	require.Equalf(t, expectedRegionDataBlock.Continent, actualRegionDataBlock.Continent, "RegionDataBlock Continent incorrect")

	// PlanDataBlock
	expectedPlanDataBlock := blocks.PlanDataBlock{VID: "vhf-8c-32gb", VCPUCount: 8, RAM: 32768, Disk: 512, DiskCount: 1, Bandwidth: 6144, MonthlyCost: 192, PlanType: "vhf", DataBlock: blocks.DataBlock{Type: blocks.PlanBlockType, Name: "vhf_32gb"}}
	actualPlanDataBlock := config.Data[blocks.PlanBlockType]["vhf_32gb"].(*blocks.PlanDataBlock)
	require.Equalf(t, blocks.PlanBlockType, actualPlanDataBlock.BlockType(), "PlanDataBlock BlockType() incorrect")
	require.Equalf(t, expectedPlanDataBlock.Name, actualPlanDataBlock.BlockName(), "PlanDataBlock BlockName() incorrect")
	require.Equalf(t, expectedPlanDataBlock.VID, actualPlanDataBlock.VID, "PlanDataBlock VID incorrect")
	require.Equalf(t, expectedPlanDataBlock.VCPUCount, actualPlanDataBlock.VCPUCount, "PlanDataBlock VCPUCount incorrect")
	require.Equalf(t, expectedPlanDataBlock.RAM, actualPlanDataBlock.RAM, "PlanDataBlock RAM incorrect")
	require.Equalf(t, expectedPlanDataBlock.Disk, actualPlanDataBlock.Disk, "PlanDataBlock Disk incorrect")

	// OSDataBlock
	expectedOSDataBlock := blocks.OSDataBlock{VID: 362, OSName: "CentOS 7 x64", Arch: "x64", Family: "centos", DataBlock: blocks.DataBlock{Type: blocks.OSBlockType, Name: "centos_7"}}
	actualOSDataBlock := config.Data[blocks.OSBlockType]["centos_7"].(*blocks.OSDataBlock)
	require.Equalf(t, blocks.OSBlockType, actualOSDataBlock.BlockType(), "OSDataBlock BlockType() incorrect")
	require.Equalf(t, expectedOSDataBlock.Name, actualOSDataBlock.BlockName(), "OSDataBlock BlockName() incorrect")
	require.Equalf(t, expectedOSDataBlock.VID, actualOSDataBlock.VID, "OSDataBlock VID incorrect")
	require.Equalf(t, expectedOSDataBlock.OSName, actualOSDataBlock.OSName, "OSDataBlock OSName incorrect")
	require.Equalf(t, expectedOSDataBlock.Arch, actualOSDataBlock.Arch, "OSDataBlock Arch incorrect")
	require.Equalf(t, expectedOSDataBlock.Family, actualOSDataBlock.Family, "OSDataBlock Family incorrect")

	// SSHKeyBlock
	expectedSSHKeyBlock := blocks.SSHKeyBlock{VID: "cb676a46-66fd-4dfb-b839-443f2e6c0b60", SSHKey: "ssh-rsa AAAAB3NzaC1yc2E", ResourceBlock: blocks.ResourceBlock{Name: "my_key", Type: blocks.SSHKeyBlockType}}
	actualSSHKeyBlock := config.Resources[blocks.SSHKeyBlockType]["my_key"].(*blocks.SSHKeyBlock)
	require.Equalf(t, blocks.SSHKeyBlockType, actualSSHKeyBlock.BlockType(), "SSHKeyBlock BlockType() incorrect")
	require.Equalf(t, expectedSSHKeyBlock.Name, actualSSHKeyBlock.BlockName(), "SSHKeyBlock BlockName() incorrect")
	require.Equalf(t, expectedSSHKeyBlock.VID, actualSSHKeyBlock.VID, "SSHKeyBlock VID incorrect")
	require.Equalf(t, expectedSSHKeyBlock.SSHKey, actualSSHKeyBlock.SSHKey, "SSHKeyBlock SSHKey incorrect")

	// StartupScriptBlock
	expectedStartupScriptBlock := blocks.StartupScriptBlock{Script: "IyEvYmluL2Jhc2gKZWNobyAnaGVsbG8gd29ybGQn", VID: "cb676a46-66fd-4dfb-b839-443f2e6c0b60", VDateCreated: "2020-10-10T01:56:20+00:00", VDateModified: "2020-10-10T01:59:20+00:00", VType: "pxe", ResourceBlock: blocks.ResourceBlock{Name: "my_script", Type: blocks.StartupScriptBlockType}}
	actualStartupScriptBlock := config.Resources[blocks.StartupScriptBlockType]["my_script"].(*blocks.StartupScriptBlock)
	require.Equalf(t, blocks.StartupScriptBlockType, actualStartupScriptBlock.BlockType(), "StartupScriptBlock BlockType() incorrect")
	require.Equalf(t, expectedStartupScriptBlock.Name, actualStartupScriptBlock.BlockName(), "StartupScriptBlock BlockName() incorrect")
	require.Equalf(t, expectedStartupScriptBlock.VID, actualStartupScriptBlock.VID, "StartupScriptBlock VID incorrect")
	require.Equalf(t, expectedStartupScriptBlock.VDateCreated, actualStartupScriptBlock.VDateCreated, "StartupScriptBlock VDateCreated incorrect")
	require.Equalf(t, expectedStartupScriptBlock.VDateModified, actualStartupScriptBlock.VDateModified, "StartupScriptBlock VDateModified incorrect")
	require.Equalf(t, expectedStartupScriptBlock.VType, actualStartupScriptBlock.VType, "StartupScriptBlock VType incorrect")
	require.Equalf(t, expectedStartupScriptBlock.Script, actualStartupScriptBlock.Script, "StartupScriptBlock Script incorrect")

	// InstanceBlock
	expectedInstanceBlock := blocks.InstanceBlock{
		VID:             "4f0f12e5-1f84-404f-aa84-85f431ea5ec2",
		OsID:            expectedOSDataBlock.VID,
		Region:          expectedRegionDataBlock.VID,
		Plan:            expectedPlanDataBlock.VID,
		SshKeyID:        expectedSSHKeyBlock.VID,
		StartupScriptID: expectedStartupScriptBlock.VID,
		Os:              "CentOS 7 x64",
		RAM:             1024,
		Disk:            25,
		VCPUCount:       1,
		Status:          "pending",
		ResourceBlock:   blocks.ResourceBlock{Name: "my_vps", Type: blocks.InstanceBlockType},
	}
	actualInstanceBlock := config.Resources[blocks.InstanceBlockType]["my_vps"].(*blocks.InstanceBlock)
	require.Equalf(t, blocks.InstanceBlockType, actualInstanceBlock.BlockType(), "InstanceBlock BlockType() incorrect")
	require.Equalf(t, expectedInstanceBlock.Name, actualInstanceBlock.BlockName(), "InstanceBlock BlockName() incorrect")
	require.Equalf(t, expectedInstanceBlock.VID, actualInstanceBlock.VID, "InstanceBlock VID incorrect")
	require.Equalf(t, expectedInstanceBlock.Os, actualInstanceBlock.Os, "InstanceBlock Os incorrect")
	require.Equalf(t, expectedInstanceBlock.OsID, actualInstanceBlock.OsID, "InstanceBlock OsID incorrect")
	require.Equalf(t, expectedInstanceBlock.Region, actualInstanceBlock.Region, "InstanceBlock Region incorrect")
	require.Equalf(t, expectedInstanceBlock.Plan, actualInstanceBlock.Plan, "InstanceBlock Plan incorrect")
	require.Equalf(t, expectedInstanceBlock.SshKeyID, actualInstanceBlock.SshKeyID, "InstanceBlock SshKeyID incorrect")
	require.Equalf(t, expectedInstanceBlock.StartupScriptID, actualInstanceBlock.StartupScriptID, "InstanceBlock StartupScriptID incorrect")
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
			if err := EvaluateConfig(GetEvalContext(), tt.args.config, tt.args.vc); (err != nil) != tt.wantErr {
				t.Errorf("EvaluateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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

	t.Fail()
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
	mockOSService.EXPECT().List(gomock.Any(), gomock.Any()).Return([]govultr.OS{
		{
			ID:     362,
			Name:   "CentOS 7 x64",
			Arch:   "x64",
			Family: "centos",
		},
	}, &govultr.Meta{}, &http.Response{}, nil).AnyTimes()
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
		Os:               "CentOS 7 x64",
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
