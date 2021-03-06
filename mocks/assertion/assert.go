// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/elethoughts-code/goasserts/assertion (interfaces: PublicTB)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPublicTB is a mock of PublicTB interface.
type MockPublicTB struct {
	ctrl     *gomock.Controller
	recorder *MockPublicTBMockRecorder
}

// MockPublicTBMockRecorder is the mock recorder for MockPublicTB.
type MockPublicTBMockRecorder struct {
	mock *MockPublicTB
}

// NewMockPublicTB creates a new mock instance.
func NewMockPublicTB(ctrl *gomock.Controller) *MockPublicTB {
	mock := &MockPublicTB{ctrl: ctrl}
	mock.recorder = &MockPublicTBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublicTB) EXPECT() *MockPublicTBMockRecorder {
	return m.recorder
}

// Cleanup mocks base method.
func (m *MockPublicTB) Cleanup(arg0 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Cleanup", arg0)
}

// Cleanup indicates an expected call of Cleanup.
func (mr *MockPublicTBMockRecorder) Cleanup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockPublicTB)(nil).Cleanup), arg0)
}

// Error mocks base method.
func (m *MockPublicTB) Error(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockPublicTBMockRecorder) Error(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockPublicTB)(nil).Error), arg0...)
}

// Errorf mocks base method.
func (m *MockPublicTB) Errorf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockPublicTBMockRecorder) Errorf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockPublicTB)(nil).Errorf), varargs...)
}

// Fail mocks base method.
func (m *MockPublicTB) Fail() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fail")
}

// Fail indicates an expected call of Fail.
func (mr *MockPublicTBMockRecorder) Fail() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockPublicTB)(nil).Fail))
}

// FailNow mocks base method.
func (m *MockPublicTB) FailNow() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FailNow")
}

// FailNow indicates an expected call of FailNow.
func (mr *MockPublicTBMockRecorder) FailNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FailNow", reflect.TypeOf((*MockPublicTB)(nil).FailNow))
}

// Failed mocks base method.
func (m *MockPublicTB) Failed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Failed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Failed indicates an expected call of Failed.
func (mr *MockPublicTBMockRecorder) Failed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Failed", reflect.TypeOf((*MockPublicTB)(nil).Failed))
}

// Fatal mocks base method.
func (m *MockPublicTB) Fatal(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockPublicTBMockRecorder) Fatal(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockPublicTB)(nil).Fatal), arg0...)
}

// Fatalf mocks base method.
func (m *MockPublicTB) Fatalf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatalf", varargs...)
}

// Fatalf indicates an expected call of Fatalf.
func (mr *MockPublicTBMockRecorder) Fatalf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatalf", reflect.TypeOf((*MockPublicTB)(nil).Fatalf), varargs...)
}

// Helper mocks base method.
func (m *MockPublicTB) Helper() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Helper")
}

// Helper indicates an expected call of Helper.
func (mr *MockPublicTBMockRecorder) Helper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Helper", reflect.TypeOf((*MockPublicTB)(nil).Helper))
}

// Log mocks base method.
func (m *MockPublicTB) Log(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Log", varargs...)
}

// Log indicates an expected call of Log.
func (mr *MockPublicTBMockRecorder) Log(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockPublicTB)(nil).Log), arg0...)
}

// Logf mocks base method.
func (m *MockPublicTB) Logf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Logf", varargs...)
}

// Logf indicates an expected call of Logf.
func (mr *MockPublicTBMockRecorder) Logf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logf", reflect.TypeOf((*MockPublicTB)(nil).Logf), varargs...)
}

// Name mocks base method.
func (m *MockPublicTB) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockPublicTBMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockPublicTB)(nil).Name))
}

// Skip mocks base method.
func (m *MockPublicTB) Skip(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Skip", varargs...)
}

// Skip indicates an expected call of Skip.
func (mr *MockPublicTBMockRecorder) Skip(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Skip", reflect.TypeOf((*MockPublicTB)(nil).Skip), arg0...)
}

// SkipNow mocks base method.
func (m *MockPublicTB) SkipNow() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SkipNow")
}

// SkipNow indicates an expected call of SkipNow.
func (mr *MockPublicTBMockRecorder) SkipNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SkipNow", reflect.TypeOf((*MockPublicTB)(nil).SkipNow))
}

// Skipf mocks base method.
func (m *MockPublicTB) Skipf(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Skipf", varargs...)
}

// Skipf indicates an expected call of Skipf.
func (mr *MockPublicTBMockRecorder) Skipf(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Skipf", reflect.TypeOf((*MockPublicTB)(nil).Skipf), varargs...)
}

// Skipped mocks base method.
func (m *MockPublicTB) Skipped() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Skipped")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Skipped indicates an expected call of Skipped.
func (mr *MockPublicTBMockRecorder) Skipped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Skipped", reflect.TypeOf((*MockPublicTB)(nil).Skipped))
}
