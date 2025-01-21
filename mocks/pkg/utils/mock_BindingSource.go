// Code generated by mockery v2.44.1. DO NOT EDIT.

package utils

import (
	utils "github.com/Alwanly/go-codebase/pkg/utils"
	mock "github.com/stretchr/testify/mock"
)

// MockBindingSource is an autogenerated mock type for the BindingSource type
type MockBindingSource struct {
	mock.Mock
}

type MockBindingSource_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBindingSource) EXPECT() *MockBindingSource_Expecter {
	return &MockBindingSource_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockBindingSource) Execute(_a0 *utils.Binder) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*utils.Binder) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBindingSource_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockBindingSource_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *utils.Binder
func (_e *MockBindingSource_Expecter) Execute(_a0 interface{}) *MockBindingSource_Execute_Call {
	return &MockBindingSource_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockBindingSource_Execute_Call) Run(run func(_a0 *utils.Binder)) *MockBindingSource_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*utils.Binder))
	})
	return _c
}

func (_c *MockBindingSource_Execute_Call) Return(_a0 error) *MockBindingSource_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBindingSource_Execute_Call) RunAndReturn(run func(*utils.Binder) error) *MockBindingSource_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockBindingSource creates a new instance of MockBindingSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBindingSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBindingSource {
	mock := &MockBindingSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
