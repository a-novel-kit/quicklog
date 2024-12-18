// Code generated by mockery v2.46.0. DO NOT EDIT.

package quicklogmocks

import mock "github.com/stretchr/testify/mock"

// MockMessage is an autogenerated mock type for the Message type
type MockMessage struct {
	mock.Mock
}

type MockMessage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessage) EXPECT() *MockMessage_Expecter {
	return &MockMessage_Expecter{mock: &_m.Mock}
}

// RenderJSON provides a mock function with given fields:
func (_m *MockMessage) RenderJSON() map[string]interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RenderJSON")
	}

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// MockMessage_RenderJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderJSON'
type MockMessage_RenderJSON_Call struct {
	*mock.Call
}

// RenderJSON is a helper method to define mock.On call
func (_e *MockMessage_Expecter) RenderJSON() *MockMessage_RenderJSON_Call {
	return &MockMessage_RenderJSON_Call{Call: _e.mock.On("RenderJSON")}
}

func (_c *MockMessage_RenderJSON_Call) Run(run func()) *MockMessage_RenderJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_RenderJSON_Call) Return(_a0 map[string]interface{}) *MockMessage_RenderJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_RenderJSON_Call) RunAndReturn(run func() map[string]interface{}) *MockMessage_RenderJSON_Call {
	_c.Call.Return(run)
	return _c
}

// RenderTerminal provides a mock function with given fields:
func (_m *MockMessage) RenderTerminal() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RenderTerminal")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockMessage_RenderTerminal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RenderTerminal'
type MockMessage_RenderTerminal_Call struct {
	*mock.Call
}

// RenderTerminal is a helper method to define mock.On call
func (_e *MockMessage_Expecter) RenderTerminal() *MockMessage_RenderTerminal_Call {
	return &MockMessage_RenderTerminal_Call{Call: _e.mock.On("RenderTerminal")}
}

func (_c *MockMessage_RenderTerminal_Call) Run(run func()) *MockMessage_RenderTerminal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockMessage_RenderTerminal_Call) Return(_a0 string) *MockMessage_RenderTerminal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessage_RenderTerminal_Call) RunAndReturn(run func() string) *MockMessage_RenderTerminal_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMessage creates a new instance of MockMessage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessage {
	mock := &MockMessage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
