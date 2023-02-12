package domain

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type NewFarmerRequest struct {
	Id        uint   `db:"id" json:"id"`
	FirstName string `db:"fname" json:"fname"`
	LastName  string `db:"lname" json:"lname"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"password"`
}

type FarmerResponse struct {
	Id        uint   `db:"id" json:"id"`
	FirstName string `db:"fname" json:"fname"`
	LastName  string `db:"lname" json:"lname"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"-"`
}

type NewMachineRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	BaseHourlyCharge uint   `json:"base_hourly_charge"`
	OwnerId          uint   `json:"owner_id"`
}

type MachineResponse struct {
	Id               uint   `db:"id" json:"id"`
	Name             string `db:"name" json:"name"`
	Description      string `db:"description" json:"description"`
	BaseHourlyCharge uint   `db:"base_hourly_charge" json:"base_hourly_charge"`
	OwnerId          uint   `db:"owner_id" json:"owner_id"`
}

type NewBookingRequest struct {
	MachineId uint   `json:"machine_id"`
	Date      string `json:"date"`
	Slots     []uint `json:"slots"`
	FarmerId  uint   `json:"farmer_id"`
}

type NewBookingResponse struct {
	InvoiceId   uint   `json:"invoice_id"`
	MachineId   uint   `json:"machine_id"`
	SlotsBooked []uint `json:"slots_booked"`
	TotalCost   uint   `json:"total_cost"`
}

type AvailabilityRequest struct {
	MachineId uint   `json:"machine_id"`
	Date      string `json:"date"`
}

type AvailabilityResponse struct {
	MachineId      uint   `json:"machine_id"`
	Date           string `json:"date"`
	SlotsAvailable []uint `json:"slots_available"`
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

type SlotResponse struct {
	SlotId    uint   `json:"slot_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
