// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"
)

// Contacts is an autogenerated mock type for the Contacts type
type Contacts struct {
	mock.Mock
}

// Resource provides a mock function with given fields: c
func (_m *Contacts) Resource(c *echo.Group) {
	_m.Called(c)
}

// NewContacts creates a new instance of Contacts. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewContacts(t interface {
	mock.TestingT
	Cleanup(func())
}) *Contacts {
	mock := &Contacts{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
