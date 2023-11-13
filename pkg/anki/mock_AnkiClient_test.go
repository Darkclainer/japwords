// Code generated by mockery v2.27.1. DO NOT EDIT.

package anki

import (
	context "context"

	ankiconnect "github.com/Darkclainer/japwords/pkg/anki/ankiconnect"

	mock "github.com/stretchr/testify/mock"
)

// MockAnkiClient is an autogenerated mock type for the AnkiClient type
type MockAnkiClient struct {
	mock.Mock
}

// AddNote provides a mock function with given fields: ctx, params, opts
func (_m *MockAnkiClient) AddNote(ctx context.Context, params *ankiconnect.AddNoteParams, opts *ankiconnect.AddNoteOptions) (int64, error) {
	ret := _m.Called(ctx, params, opts)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ankiconnect.AddNoteParams, *ankiconnect.AddNoteOptions) (int64, error)); ok {
		return rf(ctx, params, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ankiconnect.AddNoteParams, *ankiconnect.AddNoteOptions) int64); ok {
		r0 = rf(ctx, params, opts)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ankiconnect.AddNoteParams, *ankiconnect.AddNoteOptions) error); ok {
		r1 = rf(ctx, params, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDeck provides a mock function with given fields: ctx, name
func (_m *MockAnkiClient) CreateDeck(ctx context.Context, name string) (int64, error) {
	ret := _m.Called(ctx, name)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int64, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateModel provides a mock function with given fields: ctx, parameters
func (_m *MockAnkiClient) CreateModel(ctx context.Context, parameters *ankiconnect.CreateModelRequest) (int64, error) {
	ret := _m.Called(ctx, parameters)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ankiconnect.CreateModelRequest) (int64, error)); ok {
		return rf(ctx, parameters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ankiconnect.CreateModelRequest) int64); ok {
		r0 = rf(ctx, parameters)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ankiconnect.CreateModelRequest) error); ok {
		r1 = rf(ctx, parameters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeckNames provides a mock function with given fields: ctx
func (_m *MockAnkiClient) DeckNames(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModelFieldNames provides a mock function with given fields: ctx, modelName
func (_m *MockAnkiClient) ModelFieldNames(ctx context.Context, modelName string) ([]string, error) {
	ret := _m.Called(ctx, modelName)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(ctx, modelName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, modelName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, modelName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModelNames provides a mock function with given fields: ctx
func (_m *MockAnkiClient) ModelNames(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RequestPermission provides a mock function with given fields: ctx
func (_m *MockAnkiClient) RequestPermission(ctx context.Context) (*ankiconnect.RequestPermissionResponse, error) {
	ret := _m.Called(ctx)

	var r0 *ankiconnect.RequestPermissionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*ankiconnect.RequestPermissionResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *ankiconnect.RequestPermissionResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ankiconnect.RequestPermissionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockAnkiClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAnkiClient creates a new instance of MockAnkiClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAnkiClient(t mockConstructorTestingTNewMockAnkiClient) *MockAnkiClient {
	mock := &MockAnkiClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
