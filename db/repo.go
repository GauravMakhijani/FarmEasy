package db

import (
	"FarmEasy/domain"
	"context"
	"errors"

	logger "github.com/sirupsen/logrus"
)

type Storer interface {
	RegisterFarmer(context.Context, *domain.FarmerResponse) (err error)
	LoginFarmer(context.Context, string, string) (farmerId uint, err error)
	AddMachine(context.Context, *domain.MachineResponse) (err error)
	GetMachines(context.Context) (machines []domain.MachineResponse, err error)
	IsEmptySlot(context.Context, uint, uint, string) (isEmpty bool)
	AddBooking(context.Context, domain.Booking) (bookingId uint, err error)
	BookSlot(context.Context, domain.Slot) (err error)
	GetBaseCharge(context.Context, uint) (baseCharge uint, err error)
	GenrateInvoice(context.Context, domain.Invoice) (invoiceId uint, err error)
	GetBookedSlot(context.Context, uint, string) (map[uint]struct{}, error)
	GetAllBookings(context.Context, uint) (bookings []domain.BookingResponse, err error)
}

func (s *pgStore) RegisterFarmer(ctx context.Context, farmer *domain.FarmerResponse) (err error) {
	err = s.db.QueryRowContext(ctx, "INSERT INTO farmers (fname, lname, email, phone, address, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", farmer.FirstName, farmer.LastName, farmer.Email, farmer.Phone, farmer.Address, farmer.Password).Scan(&farmer.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting farmer")
		return
	}

	return
}

func (s *pgStore) LoginFarmer(ctx context.Context, email string, password string) (farmerId uint, err error) {

	err = s.db.QueryRowContext(ctx, "SELECT id FROM farmers WHERE email = $1 and password = $2", email, password).Scan(&farmerId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error incorrect email or password")
		return
	}

	return
}

func (s *pgStore) AddMachine(ctx context.Context, newMachine *domain.MachineResponse) (err error) {

	err = s.db.QueryRowContext(ctx, "INSERT INTO machines (name, description, base_hourly_charge, owner_id) VALUES ($1, $2, $3, $4) RETURNING id", newMachine.Name, newMachine.Description, newMachine.BaseHourlyCharge, newMachine.OwnerId).Scan(&newMachine.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting machine")
		return
	}

	return

}

func (s *pgStore) GetMachines(ctx context.Context) (machines []domain.MachineResponse, err error) {

	rows, err := s.db.QueryContext(ctx, "SELECT * FROM machines")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting machines")
		return
	}

	for rows.Next() {
		var machine domain.MachineResponse
		err = rows.Scan(&machine.Id, &machine.Name, &machine.Description, &machine.BaseHourlyCharge, &machine.OwnerId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning machines")
			return
		}

		machines = append(machines, machine)
	}

	return
}

func (s *pgStore) IsEmptySlot(ctx context.Context, machineId uint, slotId uint, date string) (isEmpty bool) {

	err := s.db.QueryRowContext(ctx, "SELECT slots_booked.id FROM slots_booked, bookings WHERE bookings.id = slots_booked.booking_id and  bookings.machine_id = $1 and slot_id = $2 and date = $3", machineId, slotId, date).Scan(&slotId)
	if err != nil {
		// logger.WithField("err", err.Error()).Error("Error checking slot")
		isEmpty = true
		return
	}

	isEmpty = false
	return
}

func (s *pgStore) AddBooking(ctx context.Context, booking domain.Booking) (bookingId uint, err error) {

	err = s.db.QueryRowContext(ctx, "INSERT INTO bookings (machine_id, farmer_id) VALUES ($1, $2) RETURNING id", booking.MachineId, booking.FarmerId).Scan(&bookingId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting booking")
		return
	}

	return

}

func (s *pgStore) BookSlot(ctx context.Context, slot domain.Slot) (err error) {

	_, err = s.db.ExecContext(ctx, "INSERT INTO slots_booked (booking_id, slot_id, date) VALUES ($1, $2, $3)", slot.BookingId, slot.SlotId, slot.Date)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error booking slot")
		return
	}

	return

}

func (s *pgStore) GetBaseCharge(ctx context.Context, machineId uint) (baseCharge uint, err error) {

	err = s.db.QueryRowContext(ctx, "SELECT base_hourly_charge FROM machines WHERE id = $1", machineId).Scan(&baseCharge)
	if err != nil {
		err = errors.New("error getting base charge")
		return
	}

	return

}

func (s *pgStore) GenrateInvoice(ctx context.Context, newInvoice domain.Invoice) (invoiceId uint, err error) {

	err = s.db.QueryRowContext(ctx, "INSERT INTO invoices (booking_id, date_generated, total_amount) VALUES ($1, $2, $3) RETURNING id", newInvoice.BookingId, newInvoice.DateGenrated, newInvoice.Amount).Scan(&invoiceId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error generating invoice")
		return
	}

	return

}

func (s *pgStore) GetBookedSlot(ctx context.Context, machineId uint, date string) (map[uint]struct{}, error) {

	rows, err := s.db.QueryContext(ctx, "select s.slot_id from slots_booked s , bookings b where s.booking_id = b.id and b.machine_id = $1 and s.date = $2", machineId, date)
	if err != nil {
		logger.WithField("err", err.Error()).Error("error getting booked slots")
		return nil, err
	}

	bookedSlots := map[uint]struct{}{}
	for rows.Next() {
		var id uint

		var err = rows.Scan(&id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning slots")
			return nil, err
		}

		bookedSlots[id] = struct{}{}
	}

	return bookedSlots, nil
}

func (s *pgStore) GetAllBookings(ctx context.Context, farmerId uint) (bookings []domain.BookingResponse, err error) {

	rows, err := s.db.QueryContext(ctx, "SELECT id,machine_id FROM bookings WHERE farmer_id = $1", farmerId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting bookings")
		return
	}

	for rows.Next() {
		var bookingId uint
		var machineId uint
		err = rows.Scan(&bookingId, &machineId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error scanning bookings")
			return
		}
		slotRows, err := s.db.QueryContext(ctx, "SELECT slot_id FROM slots_booked WHERE booking_id = $1", bookingId)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error getting slots for particular booking")
			return nil, err
		}

		var slots []uint
		for slotRows.Next() {
			var slotId uint
			err = slotRows.Scan(&slotId)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error scanning slots")
				break
			}
			slots = append(slots, slotId)
		}

		subBooking := domain.BookingResponse{
			BookingId:   bookingId,
			MachineId:   machineId,
			SlotsBooked: slots,
		}
		bookings = append(bookings, subBooking)
	}

	return
}
