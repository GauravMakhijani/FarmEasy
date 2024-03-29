// Code generated by mockery v2.19.0. DO NOT EDIT.

package mocks

import (
	domain "FarmEasy/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddMachine provides a mock function with given fields: _a0, _a1
func (_m *Service) AddMachine(_a0 context.Context, _a1 domain.NewMachineRequest) (domain.MachineResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.MachineResponse
	if rf, ok := ret.Get(0).(func(context.Context, domain.NewMachineRequest) domain.MachineResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.MachineResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.NewMachineRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BookMachine provides a mock function with given fields: _a0, _a1
func (_m *Service) BookMachine(_a0 context.Context, _a1 domain.NewBookingRequest) (domain.NewBookingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.NewBookingResponse
	if rf, ok := ret.Get(0).(func(context.Context, domain.NewBookingRequest) domain.NewBookingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.NewBookingResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.NewBookingRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllBookings provides a mock function with given fields: _a0, _a1
func (_m *Service) GetAllBookings(_a0 context.Context, _a1 uint) ([]domain.BookingResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []domain.BookingResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint) []domain.BookingResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.BookingResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllSlots provides a mock function with given fields: _a0
func (_m *Service) GetAllSlots(_a0 context.Context) ([]domain.SlotResponse, error) {
	ret := _m.Called(_a0)

	var r0 []domain.SlotResponse
	if rf, ok := ret.Get(0).(func(context.Context) []domain.SlotResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.SlotResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAvailability provides a mock function with given fields: _a0, _a1, _a2
func (_m *Service) GetAvailability(_a0 context.Context, _a1 uint, _a2 string) ([]uint, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []uint
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) []uint); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMachines provides a mock function with given fields: _a0
func (_m *Service) GetMachines(_a0 context.Context) ([]domain.MachineResponse, error) {
	ret := _m.Called(_a0)

	var r0 []domain.MachineResponse
	if rf, ok := ret.Get(0).(func(context.Context) []domain.MachineResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.MachineResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: _a0, _a1
func (_m *Service) Login(_a0 context.Context, _a1 domain.LoginRequest) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, domain.LoginRequest) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.LoginRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: _a0, _a1
func (_m *Service) Register(_a0 context.Context, _a1 domain.NewFarmerRequest) (domain.FarmerResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 domain.FarmerResponse
	if rf, ok := ret.Get(0).(func(context.Context, domain.NewFarmerRequest) domain.FarmerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.FarmerResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.NewFarmerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
