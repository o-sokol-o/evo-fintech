// Code generated by MockGen. DO NOT EDIT.
// Source: evo.go

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/o-sokol-o/evo-fintech/internal/domain"
)

// MockIServicesEVO is a mock of IServicesEVO interface.
type MockIServicesEVO struct {
	ctrl     *gomock.Controller
	recorder *MockIServicesEVOMockRecorder
}

// MockIServicesEVOMockRecorder is the mock recorder for MockIServicesEVO.
type MockIServicesEVOMockRecorder struct {
	mock *MockIServicesEVO
}

// NewMockIServicesEVO creates a new mock instance.
func NewMockIServicesEVO(ctrl *gomock.Controller) *MockIServicesEVO {
	mock := &MockIServicesEVO{ctrl: ctrl}
	mock.recorder = &MockIServicesEVOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIServicesEVO) EXPECT() *MockIServicesEVOMockRecorder {
	return m.recorder
}

// FetchExternTransactions mocks base method.
func (m *MockIServicesEVO) FetchExternTransactions(ctx context.Context, url string) (domain.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchExternTransactions", ctx, url)
	ret0, _ := ret[0].(domain.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchExternTransactions indicates an expected call of FetchExternTransactions.
func (mr *MockIServicesEVOMockRecorder) FetchExternTransactions(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchExternTransactions", reflect.TypeOf((*MockIServicesEVO)(nil).FetchExternTransactions), ctx, url)
}

// GetFilteredData mocks base method.
func (m *MockIServicesEVO) GetFilteredData(ctx context.Context, input domain.FilterSearchInput) ([]domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilteredData", ctx, input)
	ret0, _ := ret[0].([]domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilteredData indicates an expected call of GetFilteredData.
func (mr *MockIServicesEVOMockRecorder) GetFilteredData(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilteredData", reflect.TypeOf((*MockIServicesEVO)(nil).GetFilteredData), ctx, input)
}

// MockIServicesRemote is a mock of IServicesRemote interface.
type MockIServicesRemote struct {
	ctrl     *gomock.Controller
	recorder *MockIServicesRemoteMockRecorder
}

// MockIServicesRemoteMockRecorder is the mock recorder for MockIServicesRemote.
type MockIServicesRemoteMockRecorder struct {
	mock *MockIServicesRemote
}

// NewMockIServicesRemote creates a new mock instance.
func NewMockIServicesRemote(ctrl *gomock.Controller) *MockIServicesRemote {
	mock := &MockIServicesRemote{ctrl: ctrl}
	mock.recorder = &MockIServicesRemoteMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIServicesRemote) EXPECT() *MockIServicesRemoteMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockIServicesRemote) Get(ctx context.Context) ([]domain.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx)
	ret0, _ := ret[0].([]domain.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIServicesRemoteMockRecorder) Get(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIServicesRemote)(nil).Get), ctx)
}