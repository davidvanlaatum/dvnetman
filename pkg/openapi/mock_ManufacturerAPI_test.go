// Code generated by mockery v2.52.2. DO NOT EDIT.

package openapi

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockManufacturerAPI is an autogenerated mock type for the ManufacturerAPI type
type MockManufacturerAPI struct {
	mock.Mock
}

type MockManufacturerAPI_Expecter struct {
	mock *mock.Mock
}

func (_m *MockManufacturerAPI) EXPECT() *MockManufacturerAPI_Expecter {
	return &MockManufacturerAPI_Expecter{mock: &_m.Mock}
}

// CreateManufacturer provides a mock function with given fields: ctx, opts
func (_m *MockManufacturerAPI) CreateManufacturer(ctx context.Context, opts *CreateManufacturerOpts) (*Response, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for CreateManufacturer")
	}

	var r0 *Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *CreateManufacturerOpts) (*Response, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *CreateManufacturerOpts) *Response); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *CreateManufacturerOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManufacturerAPI_CreateManufacturer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateManufacturer'
type MockManufacturerAPI_CreateManufacturer_Call struct {
	*mock.Call
}

// CreateManufacturer is a helper method to define mock.On call
//   - ctx context.Context
//   - opts *CreateManufacturerOpts
func (_e *MockManufacturerAPI_Expecter) CreateManufacturer(ctx interface{}, opts interface{}) *MockManufacturerAPI_CreateManufacturer_Call {
	return &MockManufacturerAPI_CreateManufacturer_Call{Call: _e.mock.On("CreateManufacturer", ctx, opts)}
}

func (_c *MockManufacturerAPI_CreateManufacturer_Call) Run(run func(ctx context.Context, opts *CreateManufacturerOpts)) *MockManufacturerAPI_CreateManufacturer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*CreateManufacturerOpts))
	})
	return _c
}

func (_c *MockManufacturerAPI_CreateManufacturer_Call) Return(res *Response, err error) *MockManufacturerAPI_CreateManufacturer_Call {
	_c.Call.Return(res, err)
	return _c
}

func (_c *MockManufacturerAPI_CreateManufacturer_Call) RunAndReturn(run func(context.Context, *CreateManufacturerOpts) (*Response, error)) *MockManufacturerAPI_CreateManufacturer_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteManufacturer provides a mock function with given fields: ctx, opts
func (_m *MockManufacturerAPI) DeleteManufacturer(ctx context.Context, opts *DeleteManufacturerOpts) (*Response, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for DeleteManufacturer")
	}

	var r0 *Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteManufacturerOpts) (*Response, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *DeleteManufacturerOpts) *Response); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *DeleteManufacturerOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManufacturerAPI_DeleteManufacturer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteManufacturer'
type MockManufacturerAPI_DeleteManufacturer_Call struct {
	*mock.Call
}

// DeleteManufacturer is a helper method to define mock.On call
//   - ctx context.Context
//   - opts *DeleteManufacturerOpts
func (_e *MockManufacturerAPI_Expecter) DeleteManufacturer(ctx interface{}, opts interface{}) *MockManufacturerAPI_DeleteManufacturer_Call {
	return &MockManufacturerAPI_DeleteManufacturer_Call{Call: _e.mock.On("DeleteManufacturer", ctx, opts)}
}

func (_c *MockManufacturerAPI_DeleteManufacturer_Call) Run(run func(ctx context.Context, opts *DeleteManufacturerOpts)) *MockManufacturerAPI_DeleteManufacturer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*DeleteManufacturerOpts))
	})
	return _c
}

func (_c *MockManufacturerAPI_DeleteManufacturer_Call) Return(res *Response, err error) *MockManufacturerAPI_DeleteManufacturer_Call {
	_c.Call.Return(res, err)
	return _c
}

func (_c *MockManufacturerAPI_DeleteManufacturer_Call) RunAndReturn(run func(context.Context, *DeleteManufacturerOpts) (*Response, error)) *MockManufacturerAPI_DeleteManufacturer_Call {
	_c.Call.Return(run)
	return _c
}

// GetManufacturer provides a mock function with given fields: ctx, opts
func (_m *MockManufacturerAPI) GetManufacturer(ctx context.Context, opts *GetManufacturerOpts) (*Response, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for GetManufacturer")
	}

	var r0 *Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetManufacturerOpts) (*Response, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetManufacturerOpts) *Response); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetManufacturerOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManufacturerAPI_GetManufacturer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetManufacturer'
type MockManufacturerAPI_GetManufacturer_Call struct {
	*mock.Call
}

// GetManufacturer is a helper method to define mock.On call
//   - ctx context.Context
//   - opts *GetManufacturerOpts
func (_e *MockManufacturerAPI_Expecter) GetManufacturer(ctx interface{}, opts interface{}) *MockManufacturerAPI_GetManufacturer_Call {
	return &MockManufacturerAPI_GetManufacturer_Call{Call: _e.mock.On("GetManufacturer", ctx, opts)}
}

func (_c *MockManufacturerAPI_GetManufacturer_Call) Run(run func(ctx context.Context, opts *GetManufacturerOpts)) *MockManufacturerAPI_GetManufacturer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*GetManufacturerOpts))
	})
	return _c
}

func (_c *MockManufacturerAPI_GetManufacturer_Call) Return(res *Response, err error) *MockManufacturerAPI_GetManufacturer_Call {
	_c.Call.Return(res, err)
	return _c
}

func (_c *MockManufacturerAPI_GetManufacturer_Call) RunAndReturn(run func(context.Context, *GetManufacturerOpts) (*Response, error)) *MockManufacturerAPI_GetManufacturer_Call {
	_c.Call.Return(run)
	return _c
}

// ListManufacturers provides a mock function with given fields: ctx, opts
func (_m *MockManufacturerAPI) ListManufacturers(ctx context.Context, opts *ListManufacturersOpts) (*Response, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for ListManufacturers")
	}

	var r0 *Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ListManufacturersOpts) (*Response, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ListManufacturersOpts) *Response); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ListManufacturersOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManufacturerAPI_ListManufacturers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListManufacturers'
type MockManufacturerAPI_ListManufacturers_Call struct {
	*mock.Call
}

// ListManufacturers is a helper method to define mock.On call
//   - ctx context.Context
//   - opts *ListManufacturersOpts
func (_e *MockManufacturerAPI_Expecter) ListManufacturers(ctx interface{}, opts interface{}) *MockManufacturerAPI_ListManufacturers_Call {
	return &MockManufacturerAPI_ListManufacturers_Call{Call: _e.mock.On("ListManufacturers", ctx, opts)}
}

func (_c *MockManufacturerAPI_ListManufacturers_Call) Run(run func(ctx context.Context, opts *ListManufacturersOpts)) *MockManufacturerAPI_ListManufacturers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ListManufacturersOpts))
	})
	return _c
}

func (_c *MockManufacturerAPI_ListManufacturers_Call) Return(res *Response, err error) *MockManufacturerAPI_ListManufacturers_Call {
	_c.Call.Return(res, err)
	return _c
}

func (_c *MockManufacturerAPI_ListManufacturers_Call) RunAndReturn(run func(context.Context, *ListManufacturersOpts) (*Response, error)) *MockManufacturerAPI_ListManufacturers_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateManufacturer provides a mock function with given fields: ctx, opts
func (_m *MockManufacturerAPI) UpdateManufacturer(ctx context.Context, opts *UpdateManufacturerOpts) (*Response, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateManufacturer")
	}

	var r0 *Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *UpdateManufacturerOpts) (*Response, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *UpdateManufacturerOpts) *Response); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *UpdateManufacturerOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockManufacturerAPI_UpdateManufacturer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateManufacturer'
type MockManufacturerAPI_UpdateManufacturer_Call struct {
	*mock.Call
}

// UpdateManufacturer is a helper method to define mock.On call
//   - ctx context.Context
//   - opts *UpdateManufacturerOpts
func (_e *MockManufacturerAPI_Expecter) UpdateManufacturer(ctx interface{}, opts interface{}) *MockManufacturerAPI_UpdateManufacturer_Call {
	return &MockManufacturerAPI_UpdateManufacturer_Call{Call: _e.mock.On("UpdateManufacturer", ctx, opts)}
}

func (_c *MockManufacturerAPI_UpdateManufacturer_Call) Run(run func(ctx context.Context, opts *UpdateManufacturerOpts)) *MockManufacturerAPI_UpdateManufacturer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*UpdateManufacturerOpts))
	})
	return _c
}

func (_c *MockManufacturerAPI_UpdateManufacturer_Call) Return(res *Response, err error) *MockManufacturerAPI_UpdateManufacturer_Call {
	_c.Call.Return(res, err)
	return _c
}

func (_c *MockManufacturerAPI_UpdateManufacturer_Call) RunAndReturn(run func(context.Context, *UpdateManufacturerOpts) (*Response, error)) *MockManufacturerAPI_UpdateManufacturer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockManufacturerAPI creates a new instance of MockManufacturerAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockManufacturerAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockManufacturerAPI {
	mock := &MockManufacturerAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
