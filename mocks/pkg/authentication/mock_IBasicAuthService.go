// Code generated by mockery v2.44.1. DO NOT EDIT.

package authentication

import mock "github.com/stretchr/testify/mock"

// MockIBasicAuthService is an autogenerated mock type for the IBasicAuthService type
type MockIBasicAuthService struct {
	mock.Mock
}

type MockIBasicAuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIBasicAuthService) EXPECT() *MockIBasicAuthService_Expecter {
	return &MockIBasicAuthService_Expecter{mock: &_m.Mock}
}

// DecodeFromHeader provides a mock function with given fields: auth
func (_m *MockIBasicAuthService) DecodeFromHeader(auth string) (string, string) {
	ret := _m.Called(auth)

	if len(ret) == 0 {
		panic("no return value specified for DecodeFromHeader")
	}

	var r0 string
	var r1 string
	if rf, ok := ret.Get(0).(func(string) (string, string)); ok {
		return rf(auth)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(auth)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(auth)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// MockIBasicAuthService_DecodeFromHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DecodeFromHeader'
type MockIBasicAuthService_DecodeFromHeader_Call struct {
	*mock.Call
}

// DecodeFromHeader is a helper method to define mock.On call
//   - auth string
func (_e *MockIBasicAuthService_Expecter) DecodeFromHeader(auth interface{}) *MockIBasicAuthService_DecodeFromHeader_Call {
	return &MockIBasicAuthService_DecodeFromHeader_Call{Call: _e.mock.On("DecodeFromHeader", auth)}
}

func (_c *MockIBasicAuthService_DecodeFromHeader_Call) Run(run func(auth string)) *MockIBasicAuthService_DecodeFromHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIBasicAuthService_DecodeFromHeader_Call) Return(_a0 string, _a1 string) *MockIBasicAuthService_DecodeFromHeader_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIBasicAuthService_DecodeFromHeader_Call) RunAndReturn(run func(string) (string, string)) *MockIBasicAuthService_DecodeFromHeader_Call {
	_c.Call.Return(run)
	return _c
}

// Validate provides a mock function with given fields: username, password
func (_m *MockIBasicAuthService) Validate(username string, password string) bool {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockIBasicAuthService_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type MockIBasicAuthService_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *MockIBasicAuthService_Expecter) Validate(username interface{}, password interface{}) *MockIBasicAuthService_Validate_Call {
	return &MockIBasicAuthService_Validate_Call{Call: _e.mock.On("Validate", username, password)}
}

func (_c *MockIBasicAuthService_Validate_Call) Run(run func(username string, password string)) *MockIBasicAuthService_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockIBasicAuthService_Validate_Call) Return(_a0 bool) *MockIBasicAuthService_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIBasicAuthService_Validate_Call) RunAndReturn(run func(string, string) bool) *MockIBasicAuthService_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIBasicAuthService creates a new instance of MockIBasicAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIBasicAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIBasicAuthService {
	mock := &MockIBasicAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
