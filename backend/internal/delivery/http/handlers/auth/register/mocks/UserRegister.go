// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UserRegister is an autogenerated mock type for the UserRegister type
type UserRegister struct {
	mock.Mock
}

// Register provides a mock function with given fields: ctx, username, password
func (_m *UserRegister) Register(ctx context.Context, username string, password string) (uuid.UUID, error) {
	ret := _m.Called(ctx, username, password)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (uuid.UUID, error)); ok {
		return rf(ctx, username, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) uuid.UUID); ok {
		r0 = rf(ctx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserRegister interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserRegister creates a new instance of UserRegister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserRegister(t mockConstructorTestingTNewUserRegister) *UserRegister {
	mock := &UserRegister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}