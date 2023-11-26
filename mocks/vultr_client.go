package mocks

import (
	"context"
	"net/http"

	"github.com/vultr/govultr/v3"
	gomock "go.uber.org/mock/gomock"
	_ "go.uber.org/mock/mockgen/model"
)

// OSService is the interface to interact with the operating system endpoint on the Vultr API
// Link : https://www.vultr.com/api/#tag/os
//
//go:generate mockgen -destination=OSService.go -package=mocks github.com/bensooraj/griffon/mocks OSService
type OSService interface {
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.OS, *govultr.Meta, *http.Response, error)
}

// PlanService is the interface to interact with the Plans endpoints on the Vultr API
// Link : https://www.vultr.com/api/#tag/plans
//
//go:generate mockgen -destination=PlanService.go -package=mocks github.com/bensooraj/griffon/mocks PlanService
type PlanService interface {
	List(ctx context.Context, planType string, options *govultr.ListOptions) ([]govultr.Plan, *govultr.Meta, *http.Response, error)
	ListBareMetal(ctx context.Context, options *govultr.ListOptions) ([]govultr.BareMetalPlan, *govultr.Meta, *http.Response, error)
}

// RegionService is the interface to interact with Region endpoints on the Vultr API
// Link : https://www.vultr.com/api/#tag/region
//
//go:generate mockgen -destination=RegionService.go -package=mocks github.com/bensooraj/griffon/mocks RegionService
type RegionService interface {
	Availability(ctx context.Context, regionID string, planType string) (*govultr.PlanAvailability, *http.Response, error)
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Region, *govultr.Meta, *http.Response, error)
}

// SSHKeyService is the interface to interact with the SSH Key endpoints on the Vultr API
// Link : https://www.vultr.com/api/#tag/ssh
//
//go:generate mockgen -destination=SSHKeyService.go -package=mocks github.com/bensooraj/griffon/mocks SSHKeyService
type SSHKeyService interface { //nolint:dupl
	Create(ctx context.Context, sshKeyReq *govultr.SSHKeyReq) (*govultr.SSHKey, *http.Response, error)
	Get(ctx context.Context, sshKeyID string) (*govultr.SSHKey, *http.Response, error)
	Update(ctx context.Context, sshKeyID string, sshKeyReq *govultr.SSHKeyReq) error
	Delete(ctx context.Context, sshKeyID string) error
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.SSHKey, *govultr.Meta, *http.Response, error)
}

// StartupScriptService is the interface to interact with the startup script endpoints on the Vultr API
// Link : https://www.vultr.com/api/#tag/startup
//
//go:generate mockgen -destination=StartupScriptService.go -package=mocks github.com/bensooraj/griffon/mocks StartupScriptService
type StartupScriptService interface {
	Create(ctx context.Context, req *govultr.StartupScriptReq) (*govultr.StartupScript, *http.Response, error)
	Get(ctx context.Context, scriptID string) (*govultr.StartupScript, *http.Response, error)
	Update(ctx context.Context, scriptID string, scriptReq *govultr.StartupScriptReq) error
	Delete(ctx context.Context, scriptID string) error
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.StartupScript, *govultr.Meta, *http.Response, error)
}

// InstanceService is the interface to interact with the instance endpoints on the Vultr API
// Link: https://www.vultr.com/api/#tag/instances
//
//go:generate mockgen -destination=InstanceService.go -package=mocks github.com/bensooraj/griffon/mocks InstanceService
type InstanceService interface {
	Create(ctx context.Context, instanceReq *govultr.InstanceCreateReq) (*govultr.Instance, *http.Response, error)
	Get(ctx context.Context, instanceID string) (*govultr.Instance, *http.Response, error)
	Update(ctx context.Context, instanceID string, instanceReq *govultr.InstanceUpdateReq) (*govultr.Instance, *http.Response, error)
	Delete(ctx context.Context, instanceID string) error
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.Instance, *govultr.Meta, *http.Response, error)

	Start(ctx context.Context, instanceID string) error
	Halt(ctx context.Context, instanceID string) error
	Reboot(ctx context.Context, instanceID string) error
	Reinstall(ctx context.Context, instanceID string, reinstallReq *govultr.ReinstallReq) (*govultr.Instance, *http.Response, error)

	MassStart(ctx context.Context, instanceList []string) error
	MassHalt(ctx context.Context, instanceList []string) error
	MassReboot(ctx context.Context, instanceList []string) error

	Restore(ctx context.Context, instanceID string, restoreReq *govultr.RestoreReq) (*http.Response, error)

	GetBandwidth(ctx context.Context, instanceID string) (*govultr.Bandwidth, *http.Response, error)
	GetNeighbors(ctx context.Context, instanceID string) (*govultr.Neighbors, *http.Response, error)

	// Deprecated: ListPrivateNetworks should no longer be used. Instead, use ListVPCInfo.
	ListPrivateNetworks(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.PrivateNetwork, *govultr.Meta, *http.Response, error)
	// Deprecated: AttachPrivateNetwork should no longer be used. Instead, use AttachVPC.
	AttachPrivateNetwork(ctx context.Context, instanceID, networkID string) error
	// Deprecated: DetachPrivateNetwork should no longer be used. Instead, use DetachVPC.
	DetachPrivateNetwork(ctx context.Context, instanceID, networkID string) error

	ListVPCInfo(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.VPCInfo, *govultr.Meta, *http.Response, error)
	AttachVPC(ctx context.Context, instanceID, vpcID string) error
	DetachVPC(ctx context.Context, instanceID, vpcID string) error

	ListVPC2Info(ctx context.Context, instanceID string, options *govultr.ListOptions) ([]govultr.VPC2Info, *govultr.Meta, *http.Response, error)
	AttachVPC2(ctx context.Context, instanceID string, vpc2Req *govultr.AttachVPC2Req) error
	DetachVPC2(ctx context.Context, instanceID, vpcID string) error

	ISOStatus(ctx context.Context, instanceID string) (*govultr.Iso, *http.Response, error)
	AttachISO(ctx context.Context, instanceID, isoID string) (*http.Response, error)
	DetachISO(ctx context.Context, instanceID string) (*http.Response, error)

	GetBackupSchedule(ctx context.Context, instanceID string) (*govultr.BackupSchedule, *http.Response, error)
	SetBackupSchedule(ctx context.Context, instanceID string, backup *govultr.BackupScheduleReq) (*http.Response, error)

	CreateIPv4(ctx context.Context, instanceID string, reboot *bool) (*govultr.IPv4, *http.Response, error)
	ListIPv4(ctx context.Context, instanceID string, option *govultr.ListOptions) ([]govultr.IPv4, *govultr.Meta, *http.Response, error)
	DeleteIPv4(ctx context.Context, instanceID, ip string) error
	ListIPv6(ctx context.Context, instanceID string, option *govultr.ListOptions) ([]govultr.IPv6, *govultr.Meta, *http.Response, error)

	CreateReverseIPv6(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error
	ListReverseIPv6(ctx context.Context, instanceID string) ([]govultr.ReverseIP, *http.Response, error)
	DeleteReverseIPv6(ctx context.Context, instanceID, ip string) error

	CreateReverseIPv4(ctx context.Context, instanceID string, reverseReq *govultr.ReverseIP) error
	DefaultReverseIPv4(ctx context.Context, instanceID, ip string) error

	GetUserData(ctx context.Context, instanceID string) (*govultr.UserData, *http.Response, error)

	GetUpgrades(ctx context.Context, instanceID string) (*govultr.Upgrades, *http.Response, error)
}

func NewMockVultrClient(ctrl *gomock.Controller) *govultr.Client {
	vultrClient := govultr.NewClient(nil)
	vultrClient.SetUserAgent("mock-app")

	vultrClient.OS = NewMockOSService(ctrl)

	return vultrClient
}
