package services

import (
	"FarmEasy/domain"
	"FarmEasy/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	service Service
	repo    *mocks.Storer
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.repo = &mocks.Storer{}
	suite.service = NewFarmService(suite.repo)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	suite.repo.AssertExpectations(suite.T())
}

func (s *ServiceTestSuite) TestFarmService_Register() {
	t := s.T()

	type args struct {
		ctx    context.Context
		farmer domain.NewFarmerRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "when the farmer data is valid",
			args: args{
				ctx: context.TODO(),
				farmer: domain.NewFarmerRequest{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567890",
					Password:  "password",
					Address:   "1234, abc street, xyz city",
				},
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("RegisterFarmer", context.TODO(), mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "when repo layer returns error",
			args: args{
				ctx: context.TODO(),
				farmer: domain.NewFarmerRequest{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567890",
					Password:  "password",
					Address:   "1234, abc street, xyz city",
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("RegisterFarmer", context.TODO(), mock.Anything).Return(errors.New("mocked error")).Once()
			},
		},
		{
			name: "when the farmer phone number already exist",
			args: args{
				ctx: context.TODO(),
				farmer: domain.NewFarmerRequest{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567", ///wrong phone number
					Password:  "password",
					Address:   "1234, abc street, xyz city",
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("RegisterFarmer", context.TODO(), mock.Anything).Return(errors.New("pq: duplicate key value violates unique constraint \"farmers_email_key\"")).Once()
			},
		},
		{
			name: "when the farmer email already exist",
			args: args{
				ctx: context.TODO(),
				farmer: domain.NewFarmerRequest{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567", ///wrong phone number
					Password:  "password",
					Address:   "1234, abc street, xyz city",
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("RegisterFarmer", context.TODO(), mock.Anything).Return(errors.New("pq: duplicate key value violates unique constraint \"farmers_phone_key\"")).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.repo)
			gotAddedFarmer, err := s.service.Register(tt.args.ctx, tt.args.farmer)
			t.Log("here", gotAddedFarmer)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.IsType(t, domain.FarmerResponse{}, gotAddedFarmer)
		})
	}
}

func (s *ServiceTestSuite) TestFarmService_Login() {
	t := s.T()

	type args struct {
		ctx   context.Context
		fAuth domain.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "when login request is valid",
			args: args{
				ctx: context.TODO(),
				fAuth: domain.LoginRequest{
					Email:    "john@gmail.com",
					Password: "password",
				},
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("LoginFarmer", context.TODO(), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(uint(1), nil).Once()

			},
		},
		{
			name: "when there is repo layer error",
			args: args{
				ctx: context.TODO(),
				fAuth: domain.LoginRequest{
					Email:    "john@gmail.com",
					Password: "password",
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("LoginFarmer", context.TODO(), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(uint(1), errors.New("mocked error")).Once()
			},
		},
		{
			name: "when either username or password is incorrect",
			args: args{
				ctx: context.TODO(),
				fAuth: domain.LoginRequest{
					Email:    "john@gmail.com",
					Password: "password",
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("LoginFarmer", context.TODO(), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(uint(1), errors.New("sql: no rows in result set")).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.repo)
			gotToken, err := s.service.Login(tt.args.ctx, tt.args.fAuth)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			var token string
			assert.IsType(t, token, gotToken)
		})
	}
}

func (s *ServiceTestSuite) TestFarmService_AddMachine() {
	t := s.T()
	type args struct {
		ctx     context.Context
		machine domain.NewMachineRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				machine: domain.NewMachineRequest{
					Name:             "Sugar Cane Harvester",
					Description:      "This is a sugar cane harvester",
					BaseHourlyCharge: 1000,
				},
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("AddMachine", context.TODO(), mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				machine: domain.NewMachineRequest{
					Name:             "Sugar Cane Harvester",
					Description:      "This is a sugar cane harvester",
					BaseHourlyCharge: 1000,
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("AddMachine", context.TODO(), mock.Anything).Return(errors.New("mocked error")).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.repo)
			gotNewMachine, err := s.service.AddMachine(tt.args.ctx, tt.args.machine)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.IsType(t, domain.MachineResponse{}, gotNewMachine)
		})
	}
}

func (s *ServiceTestSuite) TestFarmService_GetMachines() {
	t := s.T()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetMachines", context.TODO()).Return([]domain.MachineResponse{}, nil).Once()
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetMachines", context.TODO()).Return([]domain.MachineResponse{}, errors.New("mocked error")).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.repo)
			gotMachines, err := s.service.GetMachines(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.IsType(t, []domain.MachineResponse{}, gotMachines)
		})
	}
}

func (s *ServiceTestSuite) TestFarmService_BookMachine() {
	t := s.T()
	type args struct {
		ctx     context.Context
		booking domain.NewBookingRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				booking: domain.NewBookingRequest{
					MachineId: 1,
					Date:      "2021-01-01",
					Slots:     []uint{1, 2, 3},
					FarmerId:  2,
				},
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("Book", context.TODO(), mock.AnythingOfType("domain.NewBookingRequest")).Return(domain.NewBookingResponse{}, nil).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range tt.args.booking.Slots {
				s.repo.On("IsEmptySlot", context.TODO(), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(true).Once()
			}

			tt.prepare(tt.args, s.repo)

			for range tt.args.booking.Slots {
				s.repo.On("BookSlot", context.TODO(), mock.AnythingOfType("domain.Slot")).Return(nil).Once()
			}

			s.repo.On("GetBaseCharge", context.TODO(), tt.args.booking.MachineId).Return(uint(1000), nil).Once()

			s.repo.On("GenrateInvoice", context.TODO(), mock.AnythingOfType("domain.Invoice")).Return(uint(1), nil).Once()

			gotInvoice, err := s.service.BookMachine(tt.args.ctx, tt.args.booking)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.IsType(t, domain.NewBookingResponse{}, gotInvoice)
		})
	}
}

func (s *ServiceTestSuite) TestFarmService_GetAvailability() {
	t := s.T()

	type args struct {
		ctx       context.Context
		machineId uint
		date      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
				date:      "2021-01-01",
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetBookedSlot", context.TODO(), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(map[uint]struct{}{}, nil).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.repo)
			gotSlotsAvailable, err := s.service.GetAvailability(tt.args.ctx, tt.args.machineId, tt.args.date)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.IsType(t, []uint{}, gotSlotsAvailable)

		})
	}
}

func (s *ServiceTestSuite) TestFarmService_GetAllBookings() {
	t := s.T()

	type args struct {
		ctx      context.Context
		farmerId uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx:      context.TODO(),
				farmerId: uint(1),
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAllBookings", context.TODO(), mock.AnythingOfType("uint")).Return([]domain.BookingResponse{}, nil).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.repo)
			gotBookings, err := s.service.GetAllBookings(tt.args.ctx, tt.args.farmerId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.IsType(t, []domain.BookingResponse{}, gotBookings)
		})
	}
}

func (s *ServiceTestSuite) TestFarmService_GetAllSlots() {
	t := s.T()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantSlots []domain.SlotResponse
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
			},
			wantSlots: []domain.SlotResponse{
				{
					SlotId:    1,
					StartTime: "10:00",
					EndTime:   "11:00",
				},
				{
					SlotId:    2,
					StartTime: "11:00",
					EndTime:   "12:00",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotSlots, err := s.service.GetAllSlots(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.IsType(t, []domain.SlotResponse{}, gotSlots)
		})
	}
}
