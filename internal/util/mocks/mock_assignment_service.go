// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/assignment/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
	business "github.com/upassed/upassed-assignment-service/internal/service/model"
)

// AssignmentService is a mock of Service interface.
type AssignmentService struct {
	ctrl     *gomock.Controller
	recorder *AssignmentServiceMockRecorder
}

// AssignmentServiceMockRecorder is the mock recorder for AssignmentService.
type AssignmentServiceMockRecorder struct {
	mock *AssignmentService
}

// NewAssignmentService creates a new mock instance.
func NewAssignmentService(ctrl *gomock.Controller) *AssignmentService {
	mock := &AssignmentService{ctrl: ctrl}
	mock.recorder = &AssignmentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *AssignmentService) EXPECT() *AssignmentServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *AssignmentService) Create(ctx context.Context, assignment *business.FormAssignment) (*business.AssignmentCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, assignment)
	ret0, _ := ret[0].(*business.AssignmentCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *AssignmentServiceMockRecorder) Create(ctx, assignment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*AssignmentService)(nil).Create), ctx, assignment)
}

// FindByFormID mocks base method.
func (m *AssignmentService) FindByFormID(ctx context.Context, formID uuid.UUID) (*business.FormAssignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByFormID", ctx, formID)
	ret0, _ := ret[0].(*business.FormAssignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByFormID indicates an expected call of FindByFormID.
func (mr *AssignmentServiceMockRecorder) FindByFormID(ctx, formID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByFormID", reflect.TypeOf((*AssignmentService)(nil).FindByFormID), ctx, formID)
}

// FindByGroupID mocks base method.
func (m *AssignmentService) FindByGroupID(ctx context.Context, groupID uuid.UUID) (*business.GroupAssignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByGroupID", ctx, groupID)
	ret0, _ := ret[0].(*business.GroupAssignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByGroupID indicates an expected call of FindByGroupID.
func (mr *AssignmentServiceMockRecorder) FindByGroupID(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByGroupID", reflect.TypeOf((*AssignmentService)(nil).FindByGroupID), ctx, groupID)
}

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// CheckDuplicates mocks base method.
func (m *Mockrepository) CheckDuplicates(ctx context.Context, assignments []*domain.Assignment) ([]*domain.Assignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicates", ctx, assignments)
	ret0, _ := ret[0].([]*domain.Assignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDuplicates indicates an expected call of CheckDuplicates.
func (mr *MockrepositoryMockRecorder) CheckDuplicates(ctx, assignments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicates", reflect.TypeOf((*Mockrepository)(nil).CheckDuplicates), ctx, assignments)
}

// FindByFormID mocks base method.
func (m *Mockrepository) FindByFormID(ctx context.Context, formID uuid.UUID) ([]*domain.Assignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByFormID", ctx, formID)
	ret0, _ := ret[0].([]*domain.Assignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByFormID indicates an expected call of FindByFormID.
func (mr *MockrepositoryMockRecorder) FindByFormID(ctx, formID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByFormID", reflect.TypeOf((*Mockrepository)(nil).FindByFormID), ctx, formID)
}

// FindByGroupID mocks base method.
func (m *Mockrepository) FindByGroupID(ctx context.Context, groupID uuid.UUID) ([]*domain.Assignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByGroupID", ctx, groupID)
	ret0, _ := ret[0].([]*domain.Assignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByGroupID indicates an expected call of FindByGroupID.
func (mr *MockrepositoryMockRecorder) FindByGroupID(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByGroupID", reflect.TypeOf((*Mockrepository)(nil).FindByGroupID), ctx, groupID)
}

// Save mocks base method.
func (m *Mockrepository) Save(ctx context.Context, assignment []*domain.Assignment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, assignment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockrepositoryMockRecorder) Save(ctx, assignment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*Mockrepository)(nil).Save), ctx, assignment)
}
