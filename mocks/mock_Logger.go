// Code generated by mockery v2.46.0. DO NOT EDIT.

package quicklogmocks

import (
	quicklog "github.com/a-novel-kit/quicklog"
	mock "github.com/stretchr/testify/mock"
)

// MockLogger is an autogenerated mock type for the Logger type
type MockLogger struct {
	mock.Mock
}

type MockLogger_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogger) EXPECT() *MockLogger_Expecter {
	return &MockLogger_Expecter{mock: &_m.Mock}
}

// Log provides a mock function with given fields: level, message
func (_m *MockLogger) Log(level quicklog.Level, message quicklog.Message) {
	_m.Called(level, message)
}

// MockLogger_Log_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Log'
type MockLogger_Log_Call struct {
	*mock.Call
}

// Log is a helper method to define mock.On call
//   - level quicklog.Level
//   - message quicklog.Message
func (_e *MockLogger_Expecter) Log(level interface{}, message interface{}) *MockLogger_Log_Call {
	return &MockLogger_Log_Call{Call: _e.mock.On("Log", level, message)}
}

func (_c *MockLogger_Log_Call) Run(run func(level quicklog.Level, message quicklog.Message)) *MockLogger_Log_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(quicklog.Level), args[1].(quicklog.Message))
	})
	return _c
}

func (_c *MockLogger_Log_Call) Return() *MockLogger_Log_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLogger_Log_Call) RunAndReturn(run func(quicklog.Level, quicklog.Message)) *MockLogger_Log_Call {
	_c.Call.Return(run)
	return _c
}

// LogAnimated provides a mock function with given fields: message
func (_m *MockLogger) LogAnimated(message quicklog.AnimatedMessage) func() {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for LogAnimated")
	}

	var r0 func()
	if rf, ok := ret.Get(0).(func(quicklog.AnimatedMessage) func()); ok {
		r0 = rf(message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func())
		}
	}

	return r0
}

// MockLogger_LogAnimated_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LogAnimated'
type MockLogger_LogAnimated_Call struct {
	*mock.Call
}

// LogAnimated is a helper method to define mock.On call
//   - message quicklog.AnimatedMessage
func (_e *MockLogger_Expecter) LogAnimated(message interface{}) *MockLogger_LogAnimated_Call {
	return &MockLogger_LogAnimated_Call{Call: _e.mock.On("LogAnimated", message)}
}

func (_c *MockLogger_LogAnimated_Call) Run(run func(message quicklog.AnimatedMessage)) *MockLogger_LogAnimated_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(quicklog.AnimatedMessage))
	})
	return _c
}

func (_c *MockLogger_LogAnimated_Call) Return(cleaner func()) *MockLogger_LogAnimated_Call {
	_c.Call.Return(cleaner)
	return _c
}

func (_c *MockLogger_LogAnimated_Call) RunAndReturn(run func(quicklog.AnimatedMessage) func()) *MockLogger_LogAnimated_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLogger creates a new instance of MockLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLogger {
	mock := &MockLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
