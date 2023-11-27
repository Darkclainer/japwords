// Code generated by mockery v2.27.1. DO NOT EDIT.

package anki

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockStatefullClient is an autogenerated mock type for the StatefullClient type
type MockStatefullClient struct {
	mock.Mock
}

// Config provides a mock function with given fields:
func (_m *MockStatefullClient) Config() *Config {
	ret := _m.Called()

	var r0 *Config
	if rf, ok := ret.Get(0).(func() *Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Config)
		}
	}

	return r0
}

// CreateDeck provides a mock function with given fields: ctx, name
func (_m *MockStatefullClient) CreateDeck(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateDefaultNoteType provides a mock function with given fields: ctx, name
func (_m *MockStatefullClient) CreateDefaultNoteType(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetState provides a mock function with given fields: ctx
func (_m *MockStatefullClient) GetState(ctx context.Context) (*State, error) {
	ret := _m.Called(ctx)

	var r0 *State
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*State, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *State); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*State)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stop provides a mock function with given fields:
func (_m *MockStatefullClient) Stop() {
	_m.Called()
}

type mockConstructorTestingTNewMockStatefullClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockStatefullClient creates a new instance of MockStatefullClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockStatefullClient(t mockConstructorTestingTNewMockStatefullClient) *MockStatefullClient {
	mock := &MockStatefullClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
