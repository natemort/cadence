// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/uber/cadence/common/persistence/serialization (interfaces: TaskSerializer)
//
// Generated by this command:
//
//	mockgen -package serialization -destination task_serializer_mock.go github.com/uber/cadence/common/persistence/serialization TaskSerializer
//

// Package serialization is a generated GoMock package.
package serialization

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	persistence "github.com/uber/cadence/common/persistence"
)

// MockTaskSerializer is a mock of TaskSerializer interface.
type MockTaskSerializer struct {
	ctrl     *gomock.Controller
	recorder *MockTaskSerializerMockRecorder
	isgomock struct{}
}

// MockTaskSerializerMockRecorder is the mock recorder for MockTaskSerializer.
type MockTaskSerializerMockRecorder struct {
	mock *MockTaskSerializer
}

// NewMockTaskSerializer creates a new mock instance.
func NewMockTaskSerializer(ctrl *gomock.Controller) *MockTaskSerializer {
	mock := &MockTaskSerializer{ctrl: ctrl}
	mock.recorder = &MockTaskSerializerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskSerializer) EXPECT() *MockTaskSerializerMockRecorder {
	return m.recorder
}

// DeserializeTask mocks base method.
func (m *MockTaskSerializer) DeserializeTask(arg0 persistence.HistoryTaskCategory, arg1 *persistence.DataBlob) (persistence.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeserializeTask", arg0, arg1)
	ret0, _ := ret[0].(persistence.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeserializeTask indicates an expected call of DeserializeTask.
func (mr *MockTaskSerializerMockRecorder) DeserializeTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeserializeTask", reflect.TypeOf((*MockTaskSerializer)(nil).DeserializeTask), arg0, arg1)
}

// SerializeTask mocks base method.
func (m *MockTaskSerializer) SerializeTask(arg0 persistence.HistoryTaskCategory, arg1 persistence.Task) (persistence.DataBlob, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SerializeTask", arg0, arg1)
	ret0, _ := ret[0].(persistence.DataBlob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SerializeTask indicates an expected call of SerializeTask.
func (mr *MockTaskSerializerMockRecorder) SerializeTask(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SerializeTask", reflect.TypeOf((*MockTaskSerializer)(nil).SerializeTask), arg0, arg1)
}
