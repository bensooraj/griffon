// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bensooraj/griffon/mocks (interfaces: InstanceService)
//
// Generated by this command:
//
//	mockgen -destination=InstanceService.go -package=mocks github.com/bensooraj/griffon/mocks InstanceService
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"

	govultr "github.com/vultr/govultr/v3"
	gomock "go.uber.org/mock/gomock"
)

// MockInstanceService is a mock of InstanceService interface.
type MockInstanceService struct {
	ctrl     *gomock.Controller
	recorder *MockInstanceServiceMockRecorder
}

// MockInstanceServiceMockRecorder is the mock recorder for MockInstanceService.
type MockInstanceServiceMockRecorder struct {
	mock *MockInstanceService
}

// NewMockInstanceService creates a new mock instance.
func NewMockInstanceService(ctrl *gomock.Controller) *MockInstanceService {
	mock := &MockInstanceService{ctrl: ctrl}
	mock.recorder = &MockInstanceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstanceService) EXPECT() *MockInstanceServiceMockRecorder {
	return m.recorder
}

// AttachISO mocks base method.
func (m *MockInstanceService) AttachISO(arg0 context.Context, arg1, arg2 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachISO", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AttachISO indicates an expected call of AttachISO.
func (mr *MockInstanceServiceMockRecorder) AttachISO(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachISO", reflect.TypeOf((*MockInstanceService)(nil).AttachISO), arg0, arg1, arg2)
}

// AttachPrivateNetwork mocks base method.
func (m *MockInstanceService) AttachPrivateNetwork(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachPrivateNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttachPrivateNetwork indicates an expected call of AttachPrivateNetwork.
func (mr *MockInstanceServiceMockRecorder) AttachPrivateNetwork(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachPrivateNetwork", reflect.TypeOf((*MockInstanceService)(nil).AttachPrivateNetwork), arg0, arg1, arg2)
}

// AttachVPC mocks base method.
func (m *MockInstanceService) AttachVPC(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachVPC", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttachVPC indicates an expected call of AttachVPC.
func (mr *MockInstanceServiceMockRecorder) AttachVPC(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachVPC", reflect.TypeOf((*MockInstanceService)(nil).AttachVPC), arg0, arg1, arg2)
}

// AttachVPC2 mocks base method.
func (m *MockInstanceService) AttachVPC2(arg0 context.Context, arg1 string, arg2 *govultr.AttachVPC2Req) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachVPC2", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AttachVPC2 indicates an expected call of AttachVPC2.
func (mr *MockInstanceServiceMockRecorder) AttachVPC2(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachVPC2", reflect.TypeOf((*MockInstanceService)(nil).AttachVPC2), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockInstanceService) Create(arg0 context.Context, arg1 *govultr.InstanceCreateReq) (*govultr.Instance, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*govultr.Instance)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockInstanceServiceMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInstanceService)(nil).Create), arg0, arg1)
}

// CreateIPv4 mocks base method.
func (m *MockInstanceService) CreateIPv4(arg0 context.Context, arg1 string, arg2 *bool) (*govultr.IPv4, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIPv4", arg0, arg1, arg2)
	ret0, _ := ret[0].(*govultr.IPv4)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateIPv4 indicates an expected call of CreateIPv4.
func (mr *MockInstanceServiceMockRecorder) CreateIPv4(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIPv4", reflect.TypeOf((*MockInstanceService)(nil).CreateIPv4), arg0, arg1, arg2)
}

// CreateReverseIPv4 mocks base method.
func (m *MockInstanceService) CreateReverseIPv4(arg0 context.Context, arg1 string, arg2 *govultr.ReverseIP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReverseIPv4", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReverseIPv4 indicates an expected call of CreateReverseIPv4.
func (mr *MockInstanceServiceMockRecorder) CreateReverseIPv4(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReverseIPv4", reflect.TypeOf((*MockInstanceService)(nil).CreateReverseIPv4), arg0, arg1, arg2)
}

// CreateReverseIPv6 mocks base method.
func (m *MockInstanceService) CreateReverseIPv6(arg0 context.Context, arg1 string, arg2 *govultr.ReverseIP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReverseIPv6", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateReverseIPv6 indicates an expected call of CreateReverseIPv6.
func (mr *MockInstanceServiceMockRecorder) CreateReverseIPv6(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReverseIPv6", reflect.TypeOf((*MockInstanceService)(nil).CreateReverseIPv6), arg0, arg1, arg2)
}

// DefaultReverseIPv4 mocks base method.
func (m *MockInstanceService) DefaultReverseIPv4(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultReverseIPv4", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DefaultReverseIPv4 indicates an expected call of DefaultReverseIPv4.
func (mr *MockInstanceServiceMockRecorder) DefaultReverseIPv4(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultReverseIPv4", reflect.TypeOf((*MockInstanceService)(nil).DefaultReverseIPv4), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockInstanceService) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockInstanceServiceMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInstanceService)(nil).Delete), arg0, arg1)
}

// DeleteIPv4 mocks base method.
func (m *MockInstanceService) DeleteIPv4(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIPv4", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIPv4 indicates an expected call of DeleteIPv4.
func (mr *MockInstanceServiceMockRecorder) DeleteIPv4(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIPv4", reflect.TypeOf((*MockInstanceService)(nil).DeleteIPv4), arg0, arg1, arg2)
}

// DeleteReverseIPv6 mocks base method.
func (m *MockInstanceService) DeleteReverseIPv6(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteReverseIPv6", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteReverseIPv6 indicates an expected call of DeleteReverseIPv6.
func (mr *MockInstanceServiceMockRecorder) DeleteReverseIPv6(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteReverseIPv6", reflect.TypeOf((*MockInstanceService)(nil).DeleteReverseIPv6), arg0, arg1, arg2)
}

// DetachISO mocks base method.
func (m *MockInstanceService) DetachISO(arg0 context.Context, arg1 string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachISO", arg0, arg1)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DetachISO indicates an expected call of DetachISO.
func (mr *MockInstanceServiceMockRecorder) DetachISO(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachISO", reflect.TypeOf((*MockInstanceService)(nil).DetachISO), arg0, arg1)
}

// DetachPrivateNetwork mocks base method.
func (m *MockInstanceService) DetachPrivateNetwork(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachPrivateNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DetachPrivateNetwork indicates an expected call of DetachPrivateNetwork.
func (mr *MockInstanceServiceMockRecorder) DetachPrivateNetwork(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachPrivateNetwork", reflect.TypeOf((*MockInstanceService)(nil).DetachPrivateNetwork), arg0, arg1, arg2)
}

// DetachVPC mocks base method.
func (m *MockInstanceService) DetachVPC(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachVPC", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DetachVPC indicates an expected call of DetachVPC.
func (mr *MockInstanceServiceMockRecorder) DetachVPC(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachVPC", reflect.TypeOf((*MockInstanceService)(nil).DetachVPC), arg0, arg1, arg2)
}

// DetachVPC2 mocks base method.
func (m *MockInstanceService) DetachVPC2(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachVPC2", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DetachVPC2 indicates an expected call of DetachVPC2.
func (mr *MockInstanceServiceMockRecorder) DetachVPC2(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachVPC2", reflect.TypeOf((*MockInstanceService)(nil).DetachVPC2), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockInstanceService) Get(arg0 context.Context, arg1 string) (*govultr.Instance, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*govultr.Instance)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockInstanceServiceMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInstanceService)(nil).Get), arg0, arg1)
}

// GetBackupSchedule mocks base method.
func (m *MockInstanceService) GetBackupSchedule(arg0 context.Context, arg1 string) (*govultr.BackupSchedule, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackupSchedule", arg0, arg1)
	ret0, _ := ret[0].(*govultr.BackupSchedule)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBackupSchedule indicates an expected call of GetBackupSchedule.
func (mr *MockInstanceServiceMockRecorder) GetBackupSchedule(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackupSchedule", reflect.TypeOf((*MockInstanceService)(nil).GetBackupSchedule), arg0, arg1)
}

// GetBandwidth mocks base method.
func (m *MockInstanceService) GetBandwidth(arg0 context.Context, arg1 string) (*govultr.Bandwidth, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBandwidth", arg0, arg1)
	ret0, _ := ret[0].(*govultr.Bandwidth)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBandwidth indicates an expected call of GetBandwidth.
func (mr *MockInstanceServiceMockRecorder) GetBandwidth(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBandwidth", reflect.TypeOf((*MockInstanceService)(nil).GetBandwidth), arg0, arg1)
}

// GetNeighbors mocks base method.
func (m *MockInstanceService) GetNeighbors(arg0 context.Context, arg1 string) (*govultr.Neighbors, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNeighbors", arg0, arg1)
	ret0, _ := ret[0].(*govultr.Neighbors)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNeighbors indicates an expected call of GetNeighbors.
func (mr *MockInstanceServiceMockRecorder) GetNeighbors(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNeighbors", reflect.TypeOf((*MockInstanceService)(nil).GetNeighbors), arg0, arg1)
}

// GetUpgrades mocks base method.
func (m *MockInstanceService) GetUpgrades(arg0 context.Context, arg1 string) (*govultr.Upgrades, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpgrades", arg0, arg1)
	ret0, _ := ret[0].(*govultr.Upgrades)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUpgrades indicates an expected call of GetUpgrades.
func (mr *MockInstanceServiceMockRecorder) GetUpgrades(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpgrades", reflect.TypeOf((*MockInstanceService)(nil).GetUpgrades), arg0, arg1)
}

// GetUserData mocks base method.
func (m *MockInstanceService) GetUserData(arg0 context.Context, arg1 string) (*govultr.UserData, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserData", arg0, arg1)
	ret0, _ := ret[0].(*govultr.UserData)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserData indicates an expected call of GetUserData.
func (mr *MockInstanceServiceMockRecorder) GetUserData(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserData", reflect.TypeOf((*MockInstanceService)(nil).GetUserData), arg0, arg1)
}

// Halt mocks base method.
func (m *MockInstanceService) Halt(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Halt", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Halt indicates an expected call of Halt.
func (mr *MockInstanceServiceMockRecorder) Halt(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Halt", reflect.TypeOf((*MockInstanceService)(nil).Halt), arg0, arg1)
}

// ISOStatus mocks base method.
func (m *MockInstanceService) ISOStatus(arg0 context.Context, arg1 string) (*govultr.Iso, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ISOStatus", arg0, arg1)
	ret0, _ := ret[0].(*govultr.Iso)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ISOStatus indicates an expected call of ISOStatus.
func (mr *MockInstanceServiceMockRecorder) ISOStatus(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ISOStatus", reflect.TypeOf((*MockInstanceService)(nil).ISOStatus), arg0, arg1)
}

// List mocks base method.
func (m *MockInstanceService) List(arg0 context.Context, arg1 *govultr.ListOptions) ([]govultr.Instance, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]govultr.Instance)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// List indicates an expected call of List.
func (mr *MockInstanceServiceMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockInstanceService)(nil).List), arg0, arg1)
}

// ListIPv4 mocks base method.
func (m *MockInstanceService) ListIPv4(arg0 context.Context, arg1 string, arg2 *govultr.ListOptions) ([]govultr.IPv4, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIPv4", arg0, arg1, arg2)
	ret0, _ := ret[0].([]govultr.IPv4)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListIPv4 indicates an expected call of ListIPv4.
func (mr *MockInstanceServiceMockRecorder) ListIPv4(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIPv4", reflect.TypeOf((*MockInstanceService)(nil).ListIPv4), arg0, arg1, arg2)
}

// ListIPv6 mocks base method.
func (m *MockInstanceService) ListIPv6(arg0 context.Context, arg1 string, arg2 *govultr.ListOptions) ([]govultr.IPv6, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIPv6", arg0, arg1, arg2)
	ret0, _ := ret[0].([]govultr.IPv6)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListIPv6 indicates an expected call of ListIPv6.
func (mr *MockInstanceServiceMockRecorder) ListIPv6(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIPv6", reflect.TypeOf((*MockInstanceService)(nil).ListIPv6), arg0, arg1, arg2)
}

// ListPrivateNetworks mocks base method.
func (m *MockInstanceService) ListPrivateNetworks(arg0 context.Context, arg1 string, arg2 *govultr.ListOptions) ([]govultr.PrivateNetwork, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPrivateNetworks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]govultr.PrivateNetwork)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListPrivateNetworks indicates an expected call of ListPrivateNetworks.
func (mr *MockInstanceServiceMockRecorder) ListPrivateNetworks(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPrivateNetworks", reflect.TypeOf((*MockInstanceService)(nil).ListPrivateNetworks), arg0, arg1, arg2)
}

// ListReverseIPv6 mocks base method.
func (m *MockInstanceService) ListReverseIPv6(arg0 context.Context, arg1 string) ([]govultr.ReverseIP, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListReverseIPv6", arg0, arg1)
	ret0, _ := ret[0].([]govultr.ReverseIP)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListReverseIPv6 indicates an expected call of ListReverseIPv6.
func (mr *MockInstanceServiceMockRecorder) ListReverseIPv6(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListReverseIPv6", reflect.TypeOf((*MockInstanceService)(nil).ListReverseIPv6), arg0, arg1)
}

// ListVPC2Info mocks base method.
func (m *MockInstanceService) ListVPC2Info(arg0 context.Context, arg1 string, arg2 *govultr.ListOptions) ([]govultr.VPC2Info, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVPC2Info", arg0, arg1, arg2)
	ret0, _ := ret[0].([]govultr.VPC2Info)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListVPC2Info indicates an expected call of ListVPC2Info.
func (mr *MockInstanceServiceMockRecorder) ListVPC2Info(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVPC2Info", reflect.TypeOf((*MockInstanceService)(nil).ListVPC2Info), arg0, arg1, arg2)
}

// ListVPCInfo mocks base method.
func (m *MockInstanceService) ListVPCInfo(arg0 context.Context, arg1 string, arg2 *govultr.ListOptions) ([]govultr.VPCInfo, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVPCInfo", arg0, arg1, arg2)
	ret0, _ := ret[0].([]govultr.VPCInfo)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListVPCInfo indicates an expected call of ListVPCInfo.
func (mr *MockInstanceServiceMockRecorder) ListVPCInfo(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVPCInfo", reflect.TypeOf((*MockInstanceService)(nil).ListVPCInfo), arg0, arg1, arg2)
}

// MassHalt mocks base method.
func (m *MockInstanceService) MassHalt(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MassHalt", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MassHalt indicates an expected call of MassHalt.
func (mr *MockInstanceServiceMockRecorder) MassHalt(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MassHalt", reflect.TypeOf((*MockInstanceService)(nil).MassHalt), arg0, arg1)
}

// MassReboot mocks base method.
func (m *MockInstanceService) MassReboot(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MassReboot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MassReboot indicates an expected call of MassReboot.
func (mr *MockInstanceServiceMockRecorder) MassReboot(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MassReboot", reflect.TypeOf((*MockInstanceService)(nil).MassReboot), arg0, arg1)
}

// MassStart mocks base method.
func (m *MockInstanceService) MassStart(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MassStart", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MassStart indicates an expected call of MassStart.
func (mr *MockInstanceServiceMockRecorder) MassStart(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MassStart", reflect.TypeOf((*MockInstanceService)(nil).MassStart), arg0, arg1)
}

// Reboot mocks base method.
func (m *MockInstanceService) Reboot(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reboot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reboot indicates an expected call of Reboot.
func (mr *MockInstanceServiceMockRecorder) Reboot(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reboot", reflect.TypeOf((*MockInstanceService)(nil).Reboot), arg0, arg1)
}

// Reinstall mocks base method.
func (m *MockInstanceService) Reinstall(arg0 context.Context, arg1 string, arg2 *govultr.ReinstallReq) (*govultr.Instance, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reinstall", arg0, arg1, arg2)
	ret0, _ := ret[0].(*govultr.Instance)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reinstall indicates an expected call of Reinstall.
func (mr *MockInstanceServiceMockRecorder) Reinstall(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reinstall", reflect.TypeOf((*MockInstanceService)(nil).Reinstall), arg0, arg1, arg2)
}

// Restore mocks base method.
func (m *MockInstanceService) Restore(arg0 context.Context, arg1 string, arg2 *govultr.RestoreReq) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Restore indicates an expected call of Restore.
func (mr *MockInstanceServiceMockRecorder) Restore(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockInstanceService)(nil).Restore), arg0, arg1, arg2)
}

// SetBackupSchedule mocks base method.
func (m *MockInstanceService) SetBackupSchedule(arg0 context.Context, arg1 string, arg2 *govultr.BackupScheduleReq) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBackupSchedule", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetBackupSchedule indicates an expected call of SetBackupSchedule.
func (mr *MockInstanceServiceMockRecorder) SetBackupSchedule(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBackupSchedule", reflect.TypeOf((*MockInstanceService)(nil).SetBackupSchedule), arg0, arg1, arg2)
}

// Start mocks base method.
func (m *MockInstanceService) Start(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockInstanceServiceMockRecorder) Start(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockInstanceService)(nil).Start), arg0, arg1)
}

// Update mocks base method.
func (m *MockInstanceService) Update(arg0 context.Context, arg1 string, arg2 *govultr.InstanceUpdateReq) (*govultr.Instance, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*govultr.Instance)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockInstanceServiceMockRecorder) Update(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockInstanceService)(nil).Update), arg0, arg1, arg2)
}