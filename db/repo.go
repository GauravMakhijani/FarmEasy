package db

import (
	"context"
	"errors"

	logger "github.com/sirupsen/logrus"
)

type Storer interface {

	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error
	RegisterFarmer(context.Context, Farmer) (addedFarmer Farmer, err error)
	LoginFarmer(context.Context, string, string) (farmerId uint, err error)
	AddMachine(context.Context, Machine) (addedMachine Machine, err error)
	GetMachines(context.Context) (machines []Machine, err error)
	IsEmptySlot(context.Context, uint, uint, string) (isEmpty bool)
	AddBooking(context.Context, Booking) (bookingId uint, err error)
	BookSlot(context.Context, Slot) (err error)
	GetBaseCharge(context.Context, uint) (baseCharge uint, err error)
	GenrateInvoice(context.Context, Invoice) (invoiceId uint, err error)
	GetBookedSlot(context.Context, uint, string) (map[uint]struct{}, error)
	GetAllBookings(context.Context, uint) (bookings []BookingResponse, err error)
}

type Farmer struct {
	Id        uint   `db:"id" json:"id"`
	FirstName string `db:"fname" json:"fname"`
	LastName  string `db:"lname" json:"lname"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"-"`
}

type Machine struct {
	Id               uint   `db:"id" json:"id"`
	Name             string `db:"name" json:"name"`
	Description      string `db:"description" json:"description"`
	BaseHourlyCharge uint   `db:"base_hourly_charge" json:"base_hourly_charge"`
	OwnerId          uint   `db:"owner_id" json:"owner_id"`
}

type Booking struct {
	Id        uint `db:"id" json:"id"`
	MachineId uint `db:"machine_id" json:"machine_id"`
	FarmerId  uint `db:"farmer_id" json:"farmer_id"`
}

type Slot struct {
	Id        uint   `db:"id" json:"id"`
	BookingId uint   `db:"booking_id" json:"booking_id"`
	SlotId    uint   `db:"slot_number" json:"slot_number"`
	Date      string `db:"date" json:"date"`
}

type Invoice struct {
	Id           uint   `db:"id" json:"id"`
	BookingId    uint   `db:"booking_id" json:"booking_id"`
	DateGenrated string `db:"date_generated" json:"date_generated"`
	Amount       uint   `db:"amount" json:"amount"`
}

type BookingResponse struct {
	BookingId   uint   `json:"booking_id"`
	MachineId   uint   `json:"machine_id"`
	SlotsBooked []uint `json:"slots_booked"`
}

func (s *pgStore) RegisterFarmer(ctx context.Context, farmer Farmer) (addedFarmer Farmer, err error) {
	err = s.db.QueryRowContext(ctx, "INSERT INTO farmers (fname, lname, email, phone, address, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", farmer.FirstName, farmer.LastName, farmer.Email, farmer.Phone, farmer.Address, farmer.Password).Scan(&farmer.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting farmer")
		return
	}

	addedFarmer = farmer

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

func (s *pgStore) AddMachine(ctx context.Context, machine Machine) (addedMachine Machine, err error) {

	err = s.db.QueryRowContext(ctx, "INSERT INTO machines (name, description, base_hourly_charge, owner_id) VALUES ($1, $2, $3, $4) RETURNING id", machine.Name, machine.Description, machine.BaseHourlyCharge, machine.OwnerId).Scan(&machine.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting machine")
		return
	}

	addedMachine = machine

	return

}

func (s *pgStore) GetMachines(ctx context.Context) (machines []Machine, err error) {

	rows, err := s.db.QueryContext(ctx, "SELECT * FROM machines")
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error getting machines")
		return
	}

	for rows.Next() {
		var machine Machine
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

func (s *pgStore) AddBooking(ctx context.Context, booking Booking) (bookingId uint, err error) {

	err = s.db.QueryRowContext(ctx, "INSERT INTO bookings (machine_id, farmer_id) VALUES ($1, $2) RETURNING id", booking.MachineId, booking.FarmerId).Scan(&bookingId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error inserting booking")
		return
	}

	return

}

func (s *pgStore) BookSlot(ctx context.Context, slot Slot) (err error) {

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

func (s *pgStore) GenrateInvoice(ctx context.Context, newInvoice Invoice) (invoiceId uint, err error) {

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

func (s *pgStore) GetAllBookings(ctx context.Context, farmerId uint) (bookings []BookingResponse, err error) {

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

		subBooking := BookingResponse{
			BookingId:   bookingId,
			MachineId:   machineId,
			SlotsBooked: slots,
		}
		bookings = append(bookings, subBooking)
	}

	return
}
