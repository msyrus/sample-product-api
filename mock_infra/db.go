// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msyrus/simple-product-inv/infra (interfaces: DB,Row)

// Package mock_infra is a generated GoMock package.
package mock_infra

import (
	gomock "github.com/golang/mock/gomock"
	infra "github.com/msyrus/simple-product-inv/infra"
	reflect "reflect"
)

// MockDB is a mock of DB interface
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// Exec mocks base method
func (m *MockDB) Exec(arg0 string, arg1 ...interface{}) error {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Exec indicates an expected call of Exec
func (mr *MockDBMockRecorder) Exec(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDB)(nil).Exec), varargs...)
}

// Query mocks base method
func (m *MockDB) Query(arg0 string, arg1 ...interface{}) (infra.Row, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(infra.Row)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockDBMockRecorder) Query(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDB)(nil).Query), varargs...)
}

// MockRow is a mock of Row interface
type MockRow struct {
	ctrl     *gomock.Controller
	recorder *MockRowMockRecorder
}

// MockRowMockRecorder is the mock recorder for MockRow
type MockRowMockRecorder struct {
	mock *MockRow
}

// NewMockRow creates a new mock instance
func NewMockRow(ctrl *gomock.Controller) *MockRow {
	mock := &MockRow{ctrl: ctrl}
	mock.recorder = &MockRowMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRow) EXPECT() *MockRowMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockRow) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockRowMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRow)(nil).Close))
}

// Next mocks base method
func (m *MockRow) Next() bool {
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next
func (mr *MockRowMockRecorder) Next() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockRow)(nil).Next))
}

// Scan mocks base method
func (m *MockRow) Scan(arg0 ...interface{}) error {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan
func (mr *MockRowMockRecorder) Scan(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockRow)(nil).Scan), arg0...)
}