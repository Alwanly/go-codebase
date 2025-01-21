// Code generated by mockery v2.44.1. DO NOT EDIT.

package database

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	zapcore "go.uber.org/zap/zapcore"
)

// MockContextFn is an autogenerated mock type for the ContextFn type
type MockContextFn struct {
	mock.Mock
}

type MockContextFn_Expecter struct {
	mock *mock.Mock
}

func (_m *MockContextFn) EXPECT() *MockContextFn_Expecter {
	return &MockContextFn_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx
func (_m *MockContextFn) Execute(ctx context.Context) []zapcore.Field {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []zapcore.Field
	if rf, ok := ret.Get(0).(func(context.Context) []zapcore.Field); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]zapcore.Field)
		}
	}

	return r0
}

// MockContextFn_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockContextFn_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockContextFn_Expecter) Execute(ctx interface{}) *MockContextFn_Execute_Call {
	return &MockContextFn_Execute_Call{Call: _e.mock.On("Execute", ctx)}
}

func (_c *MockContextFn_Execute_Call) Run(run func(ctx context.Context)) *MockContextFn_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockContextFn_Execute_Call) Return(_a0 []zapcore.Field) *MockContextFn_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockContextFn_Execute_Call) RunAndReturn(run func(context.Context) []zapcore.Field) *MockContextFn_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockContextFn creates a new instance of MockContextFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockContextFn(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockContextFn {
	mock := &MockContextFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
