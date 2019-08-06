// Copyright 2015-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-ecs-agent/agent/utils/ioutilwrapper (interfaces: IOUtil)

// Package mock_ioutilwrapper is a generated GoMock package.
package mock_ioutilwrapper

import (
	os "os"
	reflect "reflect"

	oswrapper "github.com/aws/amazon-ecs-agent/agent/utils/oswrapper"
	gomock "github.com/golang/mock/gomock"
)

// MockIOUtil is a mock of IOUtil interface
type MockIOUtil struct {
	ctrl     *gomock.Controller
	recorder *MockIOUtilMockRecorder
}

// MockIOUtilMockRecorder is the mock recorder for MockIOUtil
type MockIOUtilMockRecorder struct {
	mock *MockIOUtil
}

// NewMockIOUtil creates a new mock instance
func NewMockIOUtil(ctrl *gomock.Controller) *MockIOUtil {
	mock := &MockIOUtil{ctrl: ctrl}
	mock.recorder = &MockIOUtilMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIOUtil) EXPECT() *MockIOUtilMockRecorder {
	return m.recorder
}

// ReadFile mocks base method
func (m *MockIOUtil) ReadFile(arg0 string) ([]byte, error) {
	ret := m.ctrl.Call(m, "ReadFile", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile
func (mr *MockIOUtilMockRecorder) ReadFile(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockIOUtil)(nil).ReadFile), arg0)
}

// TempFile mocks base method
func (m *MockIOUtil) TempFile(arg0, arg1 string) (oswrapper.File, error) {
	ret := m.ctrl.Call(m, "TempFile", arg0, arg1)
	ret0, _ := ret[0].(oswrapper.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TempFile indicates an expected call of TempFile
func (mr *MockIOUtilMockRecorder) TempFile(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TempFile", reflect.TypeOf((*MockIOUtil)(nil).TempFile), arg0, arg1)
}

// WriteFile mocks base method
func (m *MockIOUtil) WriteFile(arg0 string, arg1 []byte, arg2 os.FileMode) error {
	ret := m.ctrl.Call(m, "WriteFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile
func (mr *MockIOUtilMockRecorder) WriteFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockIOUtil)(nil).WriteFile), arg0, arg1, arg2)
}
