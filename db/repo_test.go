package db

import (
	"FarmEasy/domain"
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

type DbTestSuite struct {
	suite.Suite
	repo Storer

	mock sqlxmock.Sqlmock
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (suite *DbTestSuite) SetupTest() {
	var err error
	var db *sqlx.DB
	db, suite.mock, err = sqlxmock.Newx()
	suite.Require().NoError(err)
	suite.repo = NewPgStore(db)
}

func (suite *DbTestSuite) TearDownSuite() {

	suite.mock.ExpectClose()
}

func (s *DbTestSuite) Test_pgStore_RegisterFarmer() {
	t := s.T()
	type args struct {
		ctx    context.Context
		farmer *domain.FarmerResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				farmer: &domain.FarmerResponse{
					Id:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567890",
					Address:   "1234, abc street, xyz city",
					Password:  "password",
				},
			},
			wantErr: false,
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				farmer: &domain.FarmerResponse{
					Id:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567890",
					Address:   "1234, abc street, xyz city",
					Password:  "password",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			if tt.wantErr {
				err = errors.New("mocked error")
			} else {
				err = nil
			}
			rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)

			s.mock.ExpectQuery("INSERT INTO farmers").WithArgs(tt.args.farmer.FirstName, tt.args.farmer.LastName, tt.args.farmer.Email, tt.args.farmer.Phone, tt.args.farmer.Address, tt.args.farmer.Password).WillReturnError(err).WillReturnRows(rows)

			if err := s.repo.RegisterFarmer(tt.args.ctx, tt.args.farmer); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (s *DbTestSuite) Test_pgStore_LoginFarmer() {
	t := s.T()
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx:      context.TODO(),
				email:    "john@gmail.com",
				password: "password",
			},
			wantErr: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uint(1))
				mock.ExpectQuery("SELECT id FROM farmers").WithArgs(args.email, args.password).WillReturnRows(rows)
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx:      context.TODO(),
				email:    "john@gmail.com",
				password: "password",
			},
			wantErr: true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {

				mock.ExpectQuery("SELECT id FROM farmers").WithArgs(args.email, args.password).WillReturnError(errors.New("mocked error"))

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			tt.prepare(tt.args, s.mock)
			id, err := s.repo.LoginFarmer(context.TODO(), tt.args.email, tt.args.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, uint(1), id)
			}
		})
	}
}

func (s *DbTestSuite) Test_pgStore_AddMachine() {
	t := s.T()

	type args struct {
		ctx        context.Context
		newMachine *domain.MachineResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				newMachine: &domain.MachineResponse{
					Id:               1,
					Name:             "Machine1",
					Description:      "Machine1 Description",
					BaseHourlyCharge: 1000,
					OwnerId:          1,
				},
			},
			wantErr: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uint(1))
				mock.ExpectQuery("INSERT INTO machines").WithArgs(args.newMachine.Name, args.newMachine.Description, args.newMachine.BaseHourlyCharge, args.newMachine.OwnerId).WillReturnRows(rows)
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				newMachine: &domain.MachineResponse{
					Name:             "Machine1",
					Description:      "Machine1 Description",
					BaseHourlyCharge: 1000,
					OwnerId:          1,
				},
			},
			wantErr: true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {

				mock.ExpectQuery("INSERT INTO machines").WithArgs(args.newMachine.Name, args.newMachine.Description, args.newMachine.BaseHourlyCharge, args.newMachine.OwnerId).WillReturnError(
					errors.New("mocked error"),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			if err := s.repo.AddMachine(tt.args.ctx, tt.args.newMachine); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, uint(1), tt.args.newMachine.Id)
			}

		})
	}
}

func (s *DbTestSuite) Test_pgStore_GetMachines() {
	t := s.T()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		args         args
		wantMachines []domain.MachineResponse
		wantErr      bool
		prepare      func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
			},
			wantMachines: []domain.MachineResponse{
				{
					Id:               1,
					Name:             "Machine1",
					Description:      "Machine1 Description",
					BaseHourlyCharge: 1000,
					OwnerId:          1,
				},
				{
					Id:               2,
					Name:             "Machine2",
					Description:      "Machine2 Description",
					BaseHourlyCharge: 2000,
					OwnerId:          3,
				},
			},
			wantErr: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id", "name", "description", "base_hourly_charge", "owner_id"}).
					AddRow(uint(1), "Machine1", "Machine1 Description", 1000, uint(1)).
					AddRow(uint(2), "Machine2", "Machine2 Description", 2000, uint(3))
				mock.ExpectQuery("SELECT (.+) FROM machines").WillReturnRows(rows)

			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
			},

			wantErr: true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {

				mock.ExpectQuery("SELECT (.+) FROM machines").WillReturnError(errors.New("mocked error"))

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			_, err := s.repo.GetMachines(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {

				require.NoError(t, err)

			}
		})
	}
}

func (s *DbTestSuite) Test_pgStore_IsEmptySlot() {
	t := s.T()
	type args struct {
		ctx       context.Context
		machineId uint
		slotId    uint
		date      string
	}
	tests := []struct {
		name        string
		args        args
		wantIsEmpty bool
		prepare     func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "when slot is empty",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
				slotId:    1,
				date:      "2021-01-01",
			},
			wantIsEmpty: true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id"})
				mock.ExpectQuery("SELECT slots_booked.id FROM slots_booked, bookings WHERE bookings.id = slots_booked.booking_id and bookings.machine_id = \\$1 and slot_id = \\$2 and date = \\$3").WithArgs(args.machineId, args.slotId, args.date).WillReturnRows(rows)
			},
		},
		{
			name: "when slot is not empty",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
				slotId:    1,
				date:      "2021-01-01",
			},
			wantIsEmpty: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uint(1))
				mock.ExpectQuery("SELECT slots_booked.id FROM slots_booked, bookings WHERE bookings.id = slots_booked.booking_id and bookings.machine_id = \\$1 and slot_id = \\$2 and date = \\$3").WithArgs(args.machineId, args.slotId, args.date).WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			if gotIsEmpty := s.repo.IsEmptySlot(tt.args.ctx, tt.args.machineId, tt.args.slotId, tt.args.date); tt.wantIsEmpty {
				assert.True(t, gotIsEmpty)
			} else {
				assert.False(t, gotIsEmpty)
			}
		})
	}
}

func (s *DbTestSuite) Test_pgStore_AddBooking() {
	t := s.T()
	type args struct {
		ctx     context.Context
		booking domain.Booking
	}
	tests := []struct {
		name          string
		args          args
		wantBookingId uint
		wantErr       bool
		prepare       func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				booking: domain.Booking{
					MachineId: 1,
					FarmerId:  2,
				},
			},
			wantBookingId: 1,
			wantErr:       false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uint(1))
				mock.ExpectQuery("INSERT INTO bookings").WithArgs(args.booking.MachineId, args.booking.FarmerId).WillReturnRows(rows)
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				booking: domain.Booking{
					MachineId: 1,
					FarmerId:  2,
				},
			},
			wantBookingId: 0,
			wantErr:       true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {

				mock.ExpectQuery("INSERT INTO bookings").WithArgs(args.booking.MachineId, args.booking.FarmerId).WillReturnError(errors.New("mocked error"))

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			gotBookingId, err := s.repo.AddBooking(tt.args.ctx, tt.args.booking)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantBookingId, gotBookingId)
			}
		})
	}
}

func (s *DbTestSuite) Test_pgStore_BookSlot() {
	t := s.T()
	type args struct {
		ctx  context.Context
		slot domain.Slot
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				slot: domain.Slot{
					BookingId: 1,
					SlotId:    2,
					Date:      "2021-01-01",
				},
			},
			wantErr: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO slots_booked").WithArgs(args.slot.BookingId, args.slot.SlotId, args.slot.Date).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				slot: domain.Slot{
					BookingId: 1,
					SlotId:    2,
					Date:      "2021-01-01",
				},
			},
			wantErr: true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO slots_booked").WithArgs(args.slot.BookingId, args.slot.SlotId, args.slot.Date).WillReturnError(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			if err := s.repo.BookSlot(tt.args.ctx, tt.args.slot); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (s *DbTestSuite) Test_pgStore_GetBaseCharge() {
	t := s.T()
	type args struct {
		ctx       context.Context
		machineId uint
	}
	tests := []struct {
		name           string
		args           args
		wantBaseCharge uint
		wantErr        bool
		prepare        func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
			},
			wantBaseCharge: 100,
			wantErr:        false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"base_hourly_charge"}).AddRow(uint(100))
				mock.ExpectQuery("SELECT base_hourly_charge").WithArgs(args.machineId).WillReturnRows(rows)
			},
		},

		{
			name: "negativeTest",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
			},
			wantBaseCharge: 0,
			wantErr:        true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				mock.ExpectQuery("SELECT base_hourly_charge").WithArgs(args.machineId).WillReturnError(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			gotBaseCharge, err := s.repo.GetBaseCharge(tt.args.ctx, tt.args.machineId)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantBaseCharge, gotBaseCharge)
		})
	}
}

func (s *DbTestSuite) Test_pgStore_GenrateInvoice() {
	t := s.T()
	type args struct {
		ctx        context.Context
		newInvoice domain.Invoice
	}
	tests := []struct {
		name          string
		args          args
		wantInvoiceId uint
		wantErr       bool
		prepare       func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				newInvoice: domain.Invoice{
					BookingId:    1,
					Amount:       100,
					DateGenrated: "2021-01-01",
				},
			},
			wantInvoiceId: 1,
			wantErr:       false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO invoice").WithArgs(args.newInvoice.BookingId, args.newInvoice.DateGenrated, args.newInvoice.Amount).WillReturnRows(rows)
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				newInvoice: domain.Invoice{
					BookingId:    1,
					Amount:       100,
					DateGenrated: "2021-01-01",
				},
			},
			wantInvoiceId: 0,
			wantErr:       true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO invoice").WithArgs(args.newInvoice.BookingId, args.newInvoice.Amount, args.newInvoice.DateGenrated).WillReturnError(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			gotInvoiceId, err := s.repo.GenrateInvoice(tt.args.ctx, tt.args.newInvoice)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantInvoiceId, gotInvoiceId)
		})
	}
}

func (s *DbTestSuite) Test_pgStore_GetBookedSlot() {
	t := s.T()
	type args struct {
		ctx       context.Context
		machineId uint
		date      string
	}
	tests := []struct {
		name    string
		args    args
		want    map[uint]struct{}
		wantErr bool
		prepare func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
				date:      "2021-01-01",
			},
			want: map[uint]struct{}{
				uint(1): {},
				uint(2): {},
			},
			wantErr: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"slot_id"}).AddRow(1).AddRow(2)
				mock.ExpectQuery("select s.slot_id from slots_booked s , bookings b where s.booking_id = b.id and b.machine_id = \\$1 and s.date = \\$2").WithArgs(args.machineId, args.date).WillReturnRows(rows)
			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx:       context.TODO(),
				machineId: 1,
				date:      "2021-01-01",
			},
			want:    nil,
			wantErr: true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				mock.ExpectQuery("select s.slot_id from slots_booked s , bookings b where s.booking_id = b.id and b.machine_id = \\$1 and s.date = \\$2").WithArgs(args.machineId, args.date).WillReturnError(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			got, err := s.repo.GetBookedSlot(tt.args.ctx, tt.args.machineId, tt.args.date)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s *DbTestSuite) Test_pgStore_GetAllBookings() {
	t := s.T()
	type args struct {
		ctx      context.Context
		farmerId uint
	}
	tests := []struct {
		name         string
		args         args
		wantBookings []domain.BookingResponse
		wantErr      bool
		prepare      func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx:      context.TODO(),
				farmerId: 1,
			},
			wantBookings: []domain.BookingResponse{
				{
					BookingId:   1,
					MachineId:   1,
					SlotsBooked: []uint{1, 2, 3},
				},
			},
			wantErr: false,
			prepare: func(args args, mock sqlxmock.Sqlmock) {
				rows := sqlxmock.NewRows([]string{"id", "machine_id"}).AddRow(1, 1)
				mock.ExpectQuery("SELECT id,machine_id FROM bookings WHERE farmer_id = \\$1").WithArgs(args.farmerId).WillReturnRows(rows)
				rows = sqlxmock.NewRows([]string{"slot_id"}).AddRow(1).AddRow(2).AddRow(3)
				mock.ExpectQuery("SELECT slot_id FROM slots_booked WHERE booking_id = \\$1").WithArgs(1).WillReturnRows(rows)

			},
		},
		{
			name: "negativeTest",
			args: args{
				ctx:      context.TODO(),
				farmerId: 1,
			},
			wantBookings: nil,
			wantErr:      true,
			prepare: func(args args, mock sqlxmock.Sqlmock) {

				mock.ExpectQuery("SELECT id,machine_id FROM bookings WHERE farmer_id = \\$1").WithArgs(args.farmerId).WillReturnError(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, s.mock)
			gotBookings, err := s.repo.GetAllBookings(tt.args.ctx, tt.args.farmerId)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantBookings, gotBookings)
		})
	}
}
