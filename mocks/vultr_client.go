package mocks

import (
	"context"
	"net/http"

	"github.com/vultr/govultr/v3"
	gomock "go.uber.org/mock/gomock"
	_ "go.uber.org/mock/mockgen/model"
)

//go:generate mockgen -destination=OSService.go -package=mocks github.com/bensooraj/griffon/mocks OSService
type OSService interface {
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.OS, *govultr.Meta, *http.Response, error)
}

//go:generate mockgen -destination=StartupScriptService.go -package=mocks github.com/bensooraj/griffon/mocks StartupScriptService
type StartupScriptService interface {
	Create(ctx context.Context, req *govultr.StartupScriptReq) (*govultr.StartupScript, *http.Response, error)
	Get(ctx context.Context, scriptID string) (*govultr.StartupScript, *http.Response, error)
	Update(ctx context.Context, scriptID string, scriptReq *govultr.StartupScriptReq) error
	Delete(ctx context.Context, scriptID string) error
	List(ctx context.Context, options *govultr.ListOptions) ([]govultr.StartupScript, *govultr.Meta, *http.Response, error)
}

func NewMockVultrClient(ctrl *gomock.Controller) *govultr.Client {
	vultrClient := govultr.NewClient(nil)
	vultrClient.SetUserAgent("mock-app")

	vultrClient.OS = NewMockOSService(ctrl)

	mockStartupScriptService := NewMockStartupScriptService(ctrl)
	mockStartupScriptService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&govultr.StartupScript{
		ID:           "cb676a46-66fd-4dfb-b839-443f2e6c0b60",
		DateCreated:  "2020-10-10T01:56:20+00:00",
		DateModified: "2020-10-10T01:59:20+00:00",
		Name:         "my_key",
		Type:         "pxe",
		Script:       "ssh-rsa AAAAB3NzaC1yc2E",
	}, &http.Response{}, nil)

	vultrClient.StartupScript = mockStartupScriptService

	return vultrClient
}
