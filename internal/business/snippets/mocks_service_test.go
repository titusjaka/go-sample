// Code generated by MockGen. DO NOT EDIT.
// Source: endpoints.go

// Package snippets_test is a generated GoMock package.
package snippets_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	snippets "github.com/titusjaka/go-sample/internal/business/snippets"
	service "github.com/titusjaka/go-sample/internal/infrastructure/service"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockService) Get(ctx context.Context, id uint) (snippets.Snippet, *service.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(snippets.Snippet)
	ret1, _ := ret[1].(*service.Error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockServiceMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), ctx, id)
}

// Create mocks base method
func (m *MockService) Create(ctx context.Context, snippet snippets.Snippet) (snippets.Snippet, *service.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, snippet)
	ret0, _ := ret[0].(snippets.Snippet)
	ret1, _ := ret[1].(*service.Error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockServiceMockRecorder) Create(ctx, snippet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), ctx, snippet)
}

// List mocks base method
func (m *MockService) List(ctx context.Context, limit, offset uint) ([]snippets.Snippet, service.Pagination, *service.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset)
	ret0, _ := ret[0].([]snippets.Snippet)
	ret1, _ := ret[1].(service.Pagination)
	ret2, _ := ret[2].(*service.Error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List
func (mr *MockServiceMockRecorder) List(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockService)(nil).List), ctx, limit, offset)
}

// SoftDelete mocks base method
func (m *MockService) SoftDelete(ctx context.Context, id uint) *service.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDelete", ctx, id)
	ret0, _ := ret[0].(*service.Error)
	return ret0
}

// SoftDelete indicates an expected call of SoftDelete
func (mr *MockServiceMockRecorder) SoftDelete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDelete", reflect.TypeOf((*MockService)(nil).SoftDelete), ctx, id)
}
