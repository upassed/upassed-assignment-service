// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/assignment/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/upassed/upassed-assignment-service/internal/repository/model"
)

// AssignmentRepository is a mock of Repository interface.
type AssignmentRepository struct {
	ctrl     *gomock.Controller
	recorder *AssignmentRepositoryMockRecorder
}

// AssignmentRepositoryMockRecorder is the mock recorder for AssignmentRepository.
type AssignmentRepositoryMockRecorder struct {
	mock *AssignmentRepository
}

// NewAssignmentRepository creates a new mock instance.
func NewAssignmentRepository(ctrl *gomock.Controller) *AssignmentRepository {
	mock := &AssignmentRepository{ctrl: ctrl}
	mock.recorder = &AssignmentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *AssignmentRepository) EXPECT() *AssignmentRepositoryMockRecorder {
	return m.recorder
}

// CheckDuplicates mocks base method.
func (m *AssignmentRepository) CheckDuplicates(ctx context.Context, assignments []*domain.Assignment) ([]*domain.Assignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicates", ctx, assignments)
	ret0, _ := ret[0].([]*domain.Assignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDuplicates indicates an expected call of CheckDuplicates.
func (mr *AssignmentRepositoryMockRecorder) CheckDuplicates(ctx, assignments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicates", reflect.TypeOf((*AssignmentRepository)(nil).CheckDuplicates), ctx, assignments)
}

// FindByFormID mocks base method.
func (m *AssignmentRepository) FindByFormID(ctx context.Context, formID uuid.UUID) ([]*domain.Assignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByFormID", ctx, formID)
	ret0, _ := ret[0].([]*domain.Assignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByFormID indicates an expected call of FindByFormID.
func (mr *AssignmentRepositoryMockRecorder) FindByFormID(ctx, formID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByFormID", reflect.TypeOf((*AssignmentRepository)(nil).FindByFormID), ctx, formID)
}

// FindByGroupID mocks base method.
func (m *AssignmentRepository) FindByGroupID(ctx context.Context, groupID uuid.UUID) ([]*domain.Assignment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByGroupID", ctx, groupID)
	ret0, _ := ret[0].([]*domain.Assignment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByGroupID indicates an expected call of FindByGroupID.
func (mr *AssignmentRepositoryMockRecorder) FindByGroupID(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByGroupID", reflect.TypeOf((*AssignmentRepository)(nil).FindByGroupID), ctx, groupID)
}

// Save mocks base method.
func (m *AssignmentRepository) Save(ctx context.Context, assignments []*domain.Assignment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, assignments)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *AssignmentRepositoryMockRecorder) Save(ctx, assignments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*AssignmentRepository)(nil).Save), ctx, assignments)
}
