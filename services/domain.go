package services

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type NewFarmer struct {
	Id        uint   `db:"id" json:"id"`
	FirstName string `db:"fname" json:"fname"`
	LastName  string `db:"lname" json:"lname"`
	Email     string `db:"email" json:"email"`
	Phone     string `db:"phone" json:"phone"`
	Address   string `db:"address" json:"address"`
	Password  string `db:"password" json:"password"`
}

type NewMachine struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	BaseHourlyCharge uint   `json:"base_hourly_charge"`
	OwnerId          uint   `json:"owner_id"`
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
