// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import fsnotify "github.com/fsnotify/fsnotify"

import mock "github.com/stretchr/testify/mock"

// FSNotify is an autogenerated mock type for the FSNotify type
type FSNotify struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *FSNotify) Add(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *FSNotify) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Errors provides a mock function with given fields:
func (_m *FSNotify) Errors() chan error {
	ret := _m.Called()

	var r0 chan error
	if rf, ok := ret.Get(0).(func() chan error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan error)
		}
	}

	return r0
}

// Events provides a mock function with given fields:
func (_m *FSNotify) Events() chan fsnotify.Event {
	ret := _m.Called()

	var r0 chan fsnotify.Event
	if rf, ok := ret.Get(0).(func() chan fsnotify.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan fsnotify.Event)
		}
	}

	return r0
}

// Remove provides a mock function with given fields: _a0
func (_m *FSNotify) Remove(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
