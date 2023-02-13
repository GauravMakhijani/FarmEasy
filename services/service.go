package services

import (
	"FarmEasy/constant"
	"FarmEasy/db"
	"FarmEasy/domain"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var secretKey = []byte("I'mGoingToBeAGolangDeveloper")

type Service interface {
	Register(context.Context, domain.NewFarmerRequest) (addedFarmer domain.FarmerResponse, err error)
	Login(context.Context, domain.LoginRequest) (token string, err error)
	AddMachine(context.Context, domain.NewMachineRequest) (addedMachine domain.MachineResponse, err error)
	GetMachines(context.Context) (machines []domain.MachineResponse, err error)
	BookMachine(context.Context, domain.NewBookingRequest) (domain.NewBookingResponse, error)
	GetAvailability(context.Context, uint, string) (slotsAvailable []uint, err error)
	GetAllBookings(context.Context, uint) (bookings []domain.BookingResponse, err error)
	GetAllSlots(context.Context) (slots []domain.SlotResponse, err error)
}

type FarmService struct {
	store db.Storer
}

func NewFarmService(s db.Storer) Service {
	return &FarmService{
		store: s,
	}
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

func (s *FarmService) Register(ctx context.Context, farmer domain.NewFarmerRequest) (newFarmer domain.FarmerResponse, err error) {

	newFarmer = domain.FarmerResponse{
		FirstName: farmer.FirstName,
		LastName:  farmer.LastName,
		Email:     farmer.Email,
		Phone:     farmer.Phone,
		Address:   farmer.Address,
		Password:  farmer.Password,
	}

	newFarmer.Password = Hash_password(newFarmer.Password)

	err = s.store.RegisterFarmer(ctx, &newFarmer)
	if err != nil {
		logrus.WithField("err", err.Error()).Error("error registering farmer")
		if err.Error() == "pq: duplicate key value violates unique constraint \"farmers_email_key\"" {
			err = ErrDuplicateEmail
			return
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"farmers_phone_key\"" {
			err = ErrDuplicatePhone
		}
		return

	}

	return
}

func (s *FarmService) Login(ctx context.Context, fAuth domain.LoginRequest) (token string, err error) {
	var farmerId uint
	fAuth.Password = Hash_password(fAuth.Password)
	farmerId, err = s.store.LoginFarmer(ctx, fAuth.Email, fAuth.Password)
	if err != nil {
		logrus.WithField("err", err.Error()).Error("error login farmer")
		if err.Error() == "sql: no rows in result set" {
			err = ErrUnauthorized
		}
		return
	}

	token, err = generateJWT(farmerId)
	if err != nil {
		logrus.WithField("err", err.Error()).Error("error generating fwt token for farmer")

		return
	}
	return
}
func (s *FarmService) AddMachine(ctx context.Context, machine domain.NewMachineRequest) (newMachine domain.MachineResponse, err error) {
	newMachine = domain.MachineResponse{
		Name:             machine.Name,
		Description:      machine.Description,
		BaseHourlyCharge: machine.BaseHourlyCharge,
		OwnerId:          machine.OwnerId,
	}
	err = s.store.AddMachine(ctx, &newMachine)
	return
}

func (s *FarmService) GetMachines(ctx context.Context) (machines []domain.MachineResponse, err error) {
	machines, err = s.store.GetMachines(ctx)
	return
}

func (s *FarmService) BookMachine(ctx context.Context, booking domain.NewBookingRequest) (invoice domain.NewBookingResponse, err error) {

	invoice, err = s.store.Book(ctx, booking)
	return
}

func (s *FarmService) GetAvailability(ctx context.Context, machineId uint, date string) (slotsAvailable []uint, err error) {

	bookedSlots, err := s.store.GetBookedSlot(ctx, machineId, date)
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

func (s *FarmService) GetAllBookings(ctx context.Context, farmerId uint) (bookings []domain.BookingResponse, err error) {
	bookings, err = s.store.GetAllBookings(ctx, farmerId)
	return
}

func (s *FarmService) GetAllSlots(ctx context.Context) (slots []domain.SlotResponse, err error) {
	for i := 1; i <= 24; i++ {
		slots = append(slots, domain.SlotResponse{SlotId: uint(i),
			StartTime: constant.Slots[uint(i)].StartTime,
			EndTime:   constant.Slots[uint(i)].EndTime})
	}
	return
}
