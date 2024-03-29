// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bensooraj/griffon/mocks (interfaces: PlanService)
//
// Generated by this command:
//
//	mockgen -destination=PlanService.go -package=mocks github.com/bensooraj/griffon/mocks PlanService
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

// MockPlanService is a mock of PlanService interface.
type MockPlanService struct {
	ctrl     *gomock.Controller
	recorder *MockPlanServiceMockRecorder
}

// MockPlanServiceMockRecorder is the mock recorder for MockPlanService.
type MockPlanServiceMockRecorder struct {
	mock *MockPlanService
}

// NewMockPlanService creates a new mock instance.
func NewMockPlanService(ctrl *gomock.Controller) *MockPlanService {
	mock := &MockPlanService{ctrl: ctrl}
	mock.recorder = &MockPlanServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlanService) EXPECT() *MockPlanServiceMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockPlanService) List(arg0 context.Context, arg1 string, arg2 *govultr.ListOptions) ([]govultr.Plan, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].([]govultr.Plan)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// List indicates an expected call of List.
func (mr *MockPlanServiceMockRecorder) List(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPlanService)(nil).List), arg0, arg1, arg2)
}

// ListBareMetal mocks base method.
func (m *MockPlanService) ListBareMetal(arg0 context.Context, arg1 *govultr.ListOptions) ([]govultr.BareMetalPlan, *govultr.Meta, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBareMetal", arg0, arg1)
	ret0, _ := ret[0].([]govultr.BareMetalPlan)
	ret1, _ := ret[1].(*govultr.Meta)
	ret2, _ := ret[2].(*http.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ListBareMetal indicates an expected call of ListBareMetal.
func (mr *MockPlanServiceMockRecorder) ListBareMetal(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBareMetal", reflect.TypeOf((*MockPlanService)(nil).ListBareMetal), arg0, arg1)
}
