// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	discordgo "github.com/bwmarrin/discordgo"
	ken "github.com/zekrotja/ken"

	mock "github.com/stretchr/testify/mock"
)

// IKen is an autogenerated mock type for the IKen type
type IKen struct {
	mock.Mock
}

// Components provides a mock function with given fields:
func (_m *IKen) Components() *ken.ComponentHandler {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Components")
	}

	var r0 *ken.ComponentHandler
	if rf, ok := ret.Get(0).(func() *ken.ComponentHandler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ken.ComponentHandler)
		}
	}

	return r0
}

// GetCommandInfo provides a mock function with given fields: keyTransformer
func (_m *IKen) GetCommandInfo(keyTransformer ...ken.KeyTransformerFunc) ken.CommandInfoList {
	_va := make([]interface{}, len(keyTransformer))
	for _i := range keyTransformer {
		_va[_i] = keyTransformer[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCommandInfo")
	}

	var r0 ken.CommandInfoList
	if rf, ok := ret.Get(0).(func(...ken.KeyTransformerFunc) ken.CommandInfoList); ok {
		r0 = rf(keyTransformer...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ken.CommandInfoList)
		}
	}

	return r0
}

// RegisterCommands provides a mock function with given fields: cmds
func (_m *IKen) RegisterCommands(cmds ...ken.Command) error {
	_va := make([]interface{}, len(cmds))
	for _i := range cmds {
		_va[_i] = cmds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RegisterCommands")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...ken.Command) error); ok {
		r0 = rf(cmds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterMiddlewares provides a mock function with given fields: mws
func (_m *IKen) RegisterMiddlewares(mws ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, mws...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RegisterMiddlewares")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(mws...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Session provides a mock function with given fields:
func (_m *IKen) Session() *discordgo.Session {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Session")
	}

	var r0 *discordgo.Session
	if rf, ok := ret.Get(0).(func() *discordgo.Session); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discordgo.Session)
		}
	}

	return r0
}

// Unregister provides a mock function with given fields:
func (_m *IKen) Unregister() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Unregister")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIKen creates a new instance of IKen. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIKen(t interface {
	mock.TestingT
	Cleanup(func())
}) *IKen {
	mock := &IKen{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
