package services

import (
	"FarmEasy/db"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var secretKey = []byte("I'mGoingToBeAGolangDeveloper")

type Service interface {
	Register(context.Context, NewFarmer) (addedFarmer db.Farmer, err error)
	Login(context.Context, NewLogin) (token string, err error)
	AddMachine(context.Context, NewMachine) (addedMachine db.Machine, err error)
	GetMachines(context.Context) (machines []db.Machine, err error)
	BookMachine(context.Context, NewBooking) (db.Invoice, error)
	GetAvailability(context.Context, uint) (slotsAvailable []uint, err error)
	GetAllBookings(context.Context, uint) (bookings []db.Booking, err error)
}

type FarmService struct {
	store db.Storer
}

func NewFarmService(s db.Storer) Service {
	return &FarmService{
		store: s,
	}
}

func (s *FarmService) Register(ctx context.Context, farmer NewFarmer) (addedFarmer db.Farmer, err error) {

	newFarmer := db.Farmer{
		FirstName: farmer.FirstName,
		LastName:  farmer.LastName,
		Email:     farmer.Email,
		Phone:     farmer.Phone,
		Address:   farmer.Address,
		Password:  farmer.Password,
	}

	newFarmer.Password = Hash_password(newFarmer.Password)

	addedFarmer, err = s.store.RegisterFarmer(ctx, newFarmer)
	return
}
func generateJWT(farmerId uint) (token string, err error) {
	tokenExpirationTime := time.Now().Add(time.Hour * 24)
	tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": farmerId,
		"exp":     tokenExpirationTime.Unix(),
	})
	token, err = tokenObject.SignedString(secretKey)
	return
}

func Hash_password(password string) (hash string) {

	hsha := sha256.New()
	hsha.Write([]byte(password))
	hash = base64.URLEncoding.EncodeToString(hsha.Sum(nil))
	logrus.Info(password, " -> ", hash)
	return
}

func (s *FarmService) Login(ctx context.Context, fAuth NewLogin) (token string, err error) {
	var farmerId uint
	fAuth.Password = Hash_password(fAuth.Password)
	farmerId, err = s.store.LoginFarmer(ctx, fAuth.Email, fAuth.Password)
	if err != nil {
		return
	}
	token, err = generateJWT(farmerId)
	if err != nil {
		return
	}
	return
}
func (s *FarmService) AddMachine(ctx context.Context, machine NewMachine) (addedMachine db.Machine, err error) {
	newMachine := db.Machine{
		Name:             machine.Name,
		Description:      machine.Description,
		BaseHourlyCharge: machine.BaseHourlyCharge,
		OwnerId:          machine.OwnerId,
	}
	addedMachine, err = s.store.AddMachine(ctx, newMachine)
	return
}

func (s *FarmService) GetMachines(ctx context.Context) (machines []db.Machine, err error) {
	machines, err = s.store.GetMachines(ctx)
	return
}

func (s *FarmService) BookMachine(ctx context.Context, booking NewBooking) (invoice db.Invoice, err error) {

	for _, slot := range booking.Slots {
		empty := s.store.IsEmptySlot(ctx, booking.MachineId, slot, booking.Date)
		if !empty {

			err = errors.New("slot not empty")
			return
		}
	}

	newBooking := db.Booking{
		MachineId: booking.MachineId,
		FarmerId:  booking.FarmerId,
	}
	newBooking.Id, err = s.store.AddBooking(ctx, newBooking)
	if err != nil {
		return
	}
	for _, slot := range booking.Slots {
		newSlot := db.Slot{
			BookingId: newBooking.Id,
			SlotId:    slot,
			Date:      booking.Date,
		}
		err = s.store.BookSlot(ctx, newSlot)
		if err != nil {
			return
		}
	}
	baseCharge, err := s.store.GetBaseCharge(ctx, booking.MachineId)
	if err != nil {
		return
	}
	totalAmount := uint(len(booking.Slots)) * baseCharge
	newInvoice := db.Invoice{
		BookingId:    newBooking.Id,
		DateGenrated: time.Now().Format("2006-01-02"),
		Amount:       totalAmount,
	}
	newInvoice.Id, err = s.store.GenrateInvoice(ctx, newInvoice)

	invoice = newInvoice

	return
}

func (s *FarmService) GetAvailability(ctx context.Context, machineId uint) (slotsAvailable []uint, err error) {

	bookedSlots, err := s.store.GetBookedSlot(ctx, machineId)
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 24; i++ {
		if _, ok := bookedSlots[uint(i)]; !ok {
			slotsAvailable = append(slotsAvailable, uint(i))
		}
	}
	return
}

func (s *FarmService) GetAllBookings(ctx context.Context, farmerId uint) (bookings []db.Booking, err error) {
	bookings, err = s.store.GetAllBookings(ctx, farmerId)
	return
}
