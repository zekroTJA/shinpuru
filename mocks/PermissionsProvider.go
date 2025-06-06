// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	fiber "github.com/gofiber/fiber/v2"
	discordutil "github.com/zekroTJA/shinpuru/pkg/discordutil"

	ken "github.com/zekrotja/ken"

	mock "github.com/stretchr/testify/mock"

	pkgpermissions "github.com/zekroTJA/shinpuru/pkg/permissions"
)

// PermissionsProvider is an autogenerated mock type for the Provider type
type PermissionsProvider struct {
	mock.Mock
}

// Before provides a mock function with given fields: ctx
func (_m *PermissionsProvider) Before(ctx *ken.Ctx) (bool, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Before")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*ken.Ctx) (bool, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(*ken.Ctx) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*ken.Ctx) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckPermissions provides a mock function with given fields: s, guildID, userID, dns
func (_m *PermissionsProvider) CheckPermissions(s discordutil.ISession, guildID string, userID string, dns ...string) (bool, bool, error) {
	_va := make([]interface{}, len(dns))
	for _i := range dns {
		_va[_i] = dns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, s, guildID, userID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CheckPermissions")
	}

	var r0 bool
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string, string, ...string) (bool, bool, error)); ok {
		return rf(s, guildID, userID, dns...)
	}
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string, string, ...string) bool); ok {
		r0 = rf(s, guildID, userID, dns...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(discordutil.ISession, string, string, ...string) bool); ok {
		r1 = rf(s, guildID, userID, dns...)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(discordutil.ISession, string, string, ...string) error); ok {
		r2 = rf(s, guildID, userID, dns...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CheckSubPerm provides a mock function with given fields: ctx, subDN, explicit, message
func (_m *PermissionsProvider) CheckSubPerm(ctx ken.Context, subDN string, explicit bool, message ...string) (bool, error) {
	_va := make([]interface{}, len(message))
	for _i := range message {
		_va[_i] = message[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, subDN, explicit)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CheckSubPerm")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(ken.Context, string, bool, ...string) (bool, error)); ok {
		return rf(ctx, subDN, explicit, message...)
	}
	if rf, ok := ret.Get(0).(func(ken.Context, string, bool, ...string) bool); ok {
		r0 = rf(ctx, subDN, explicit, message...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(ken.Context, string, bool, ...string) error); ok {
		r1 = rf(ctx, subDN, explicit, message...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberPermission provides a mock function with given fields: s, guildID, memberID
func (_m *PermissionsProvider) GetMemberPermission(s discordutil.ISession, guildID string, memberID string) (pkgpermissions.PermissionArray, error) {
	ret := _m.Called(s, guildID, memberID)

	if len(ret) == 0 {
		panic("no return value specified for GetMemberPermission")
	}

	var r0 pkgpermissions.PermissionArray
	var r1 error
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string, string) (pkgpermissions.PermissionArray, error)); ok {
		return rf(s, guildID, memberID)
	}
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string, string) pkgpermissions.PermissionArray); ok {
		r0 = rf(s, guildID, memberID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pkgpermissions.PermissionArray)
		}
	}

	if rf, ok := ret.Get(1).(func(discordutil.ISession, string, string) error); ok {
		r1 = rf(s, guildID, memberID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPermissions provides a mock function with given fields: s, guildID, userID
func (_m *PermissionsProvider) GetPermissions(s discordutil.ISession, guildID string, userID string) (pkgpermissions.PermissionArray, bool, error) {
	ret := _m.Called(s, guildID, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetPermissions")
	}

	var r0 pkgpermissions.PermissionArray
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string, string) (pkgpermissions.PermissionArray, bool, error)); ok {
		return rf(s, guildID, userID)
	}
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string, string) pkgpermissions.PermissionArray); ok {
		r0 = rf(s, guildID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pkgpermissions.PermissionArray)
		}
	}

	if rf, ok := ret.Get(1).(func(discordutil.ISession, string, string) bool); ok {
		r1 = rf(s, guildID, userID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(discordutil.ISession, string, string) error); ok {
		r2 = rf(s, guildID, userID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HandleWs provides a mock function with given fields: s, required
func (_m *PermissionsProvider) HandleWs(s discordutil.ISession, required string) func(*fiber.Ctx) error {
	ret := _m.Called(s, required)

	if len(ret) == 0 {
		panic("no return value specified for HandleWs")
	}

	var r0 func(*fiber.Ctx) error
	if rf, ok := ret.Get(0).(func(discordutil.ISession, string) func(*fiber.Ctx) error); ok {
		r0 = rf(s, required)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(*fiber.Ctx) error)
		}
	}

	return r0
}

// NewPermissionsProvider creates a new instance of PermissionsProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPermissionsProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *PermissionsProvider {
	mock := &PermissionsProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
