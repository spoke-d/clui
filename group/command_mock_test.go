// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/spoke-d/clui/group (interfaces: Command)

// Package group is a generated GoMock package.
package group

import (
	gomock "github.com/golang/mock/gomock"
	flagset "github.com/spoke-d/clui/flagset"
	group "github.com/spoke-d/task/group"
	reflect "reflect"
)

// MockCommand is a mock of Command interface
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// FlagSet mocks base method
func (m *MockCommand) FlagSet() *flagset.FlagSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlagSet")
	ret0, _ := ret[0].(*flagset.FlagSet)
	return ret0
}

// FlagSet indicates an expected call of FlagSet
func (mr *MockCommandMockRecorder) FlagSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlagSet", reflect.TypeOf((*MockCommand)(nil).FlagSet))
}

// Help mocks base method
func (m *MockCommand) Help() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Help")
	ret0, _ := ret[0].(string)
	return ret0
}

// Help indicates an expected call of Help
func (mr *MockCommandMockRecorder) Help() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Help", reflect.TypeOf((*MockCommand)(nil).Help))
}

// Run mocks base method
func (m *MockCommand) Run(arg0 *group.Group) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run
func (mr *MockCommandMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCommand)(nil).Run), arg0)
}

// Synopsis mocks base method
func (m *MockCommand) Synopsis() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Synopsis")
	ret0, _ := ret[0].(string)
	return ret0
}

// Synopsis indicates an expected call of Synopsis
func (mr *MockCommandMockRecorder) Synopsis() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Synopsis", reflect.TypeOf((*MockCommand)(nil).Synopsis))
}
