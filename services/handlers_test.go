package services

import (
	"FarmEasy/api"
	"FarmEasy/domain"
	"FarmEasy/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	service *mocks.Service
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.service = &mocks.Service{}
}

func (suite *HandlerTestSuite) TearDownTest() {
	suite.service.AssertExpectations(suite.T())
}

func (s *HandlerTestSuite) Test_registerHandler() {
	t := s.T()

	t.Run("when valid register request is made", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"fname": "John", "lname": "Doe", "email": "john@gmail.com" , "phone": "1234567890", "password": "password"}`)
		r := httptest.NewRequest(http.MethodPost, "/register", (bodyReader))
		w := httptest.NewRecorder()
		ctx := r.Context()
		respBody := domain.FarmerResponse{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@gmail.com",
			Phone:     "1234567890",
			Password:  "password",
		}
		requestBody := domain.NewFarmerRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@gmail.com",
			Phone:     "1234567890",
			Password:  "password",
		}
		s.service.On("Register", ctx, requestBody).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := registerHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when error in registering user", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"fname": "John", "lname": "Doe", "email": "john@gmail.com" , "phone": "1234567890", "password": "password"}`)
		r := httptest.NewRequest(http.MethodPost, "/register", (bodyReader))
		w := httptest.NewRecorder()
		ctx := r.Context()
		respBody := api.Message{
			Msg: "Err - mocked error",
		}
		requestBody := domain.NewFarmerRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@gmail.com",
			Phone:     "1234567890",
			Password:  "password",
		}
		s.service.On("Register", ctx, requestBody).Return(domain.FarmerResponse{}, errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := registerHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when invalid register request is made, invalid email", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"fname": "John", "lname": "Doe", "email": "john@gmail" , "phone": "1234567890", "password": "password"}`)
		r := httptest.NewRequest(http.MethodPost, "/register", (bodyReader))
		w := httptest.NewRecorder()
		respBody := api.Message{
			Msg: "invalid email",
		}

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := registerHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when invalid register request is made, invalid phone", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"fname": "John", "lname": "Doe", "email": "john@gmail.com" , "phone": "123456789", "password": "password"}`)
		r := httptest.NewRequest(http.MethodPost, "/register", (bodyReader))
		w := httptest.NewRecorder()
		respBody := api.Message{
			Msg: "invalid phone number",
		}

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := registerHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(exp), w.Body.String())
	})

}

func (s *HandlerTestSuite) Test_loginHandler() {
	t := s.T()

	t.Run("when valid login request is made", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"email": "john@gmail.com" , "password": "password"}`)
		r := httptest.NewRequest(http.MethodPost, "/login", (bodyReader))
		w := httptest.NewRecorder()
		ctx := r.Context()
		respBody := domain.LoginResponse{
			Message: "Login Successful",
			Token:   "token",
		}
		requestBody := domain.LoginRequest{
			Email:    "john@gmail.com",
			Password: "password",
		}
		s.service.On("Login", ctx, requestBody).Return("token", nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := loginHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when error in logging user", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"email": "john@gmail.com" , "password": "password"}`)
		r := httptest.NewRequest(http.MethodPost, "/login", (bodyReader))
		w := httptest.NewRecorder()
		ctx := r.Context()
		respBody := api.Message{
			Msg: "mocked error",
		}
		requestBody := domain.LoginRequest{
			Email:    "john@gmail.com",
			Password: "password",
		}
		s.service.On("Login", ctx, requestBody).Return("token", errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := loginHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when invalid register request is made, invalid email", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"email": "john@gmail" , "password": "password"}`)

		r := httptest.NewRequest(http.MethodPost, "/login", (bodyReader))
		w := httptest.NewRecorder()
		respBody := api.Message{
			Msg: "invalid email",
		}

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := loginHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(exp), w.Body.String())
	})

}

func (s *HandlerTestSuite) Test_addMachineHandler() {
	t := s.T()
	t.Run("when valid add machine request is made", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"name" : "machine1", "description" : "machine1 description", "base_hourly_charge": 500, "owner_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/machines", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := domain.MachineResponse{
			Id:               1,
			Name:             "machine1",
			Description:      "machine1 description",
			BaseHourlyCharge: 500,
			OwnerId:          1,
		}
		requestBody := domain.NewMachineRequest{

			Name:             "machine1",
			Description:      "machine1 description",
			BaseHourlyCharge: 500,
			OwnerId:          1,
		}
		s.service.On("AddMachine", ctx, requestBody).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := addMachineHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when error in adding machine", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"name" : "machine1", "description" : "machine1 description", "base_hourly_charge": 500, "owner_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/machines", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := api.Message{
			Msg: "mocked error",
		}
		requestBody := domain.NewMachineRequest{

			Name:             "machine1",
			Description:      "machine1 description",
			BaseHourlyCharge: 500,
			OwnerId:          1,
		}
		s.service.On("AddMachine", ctx, requestBody).Return(domain.MachineResponse{}, errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := addMachineHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(exp), w.Body.String())
	})

}

func (s *HandlerTestSuite) Test_getMachineHandler() {
	t := s.T()
	t.Run("when valid get machine request is made", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/machines", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := []domain.MachineResponse{
			{
				Id:               1,
				Name:             "machine1",
				Description:      "machine1 description",
				BaseHourlyCharge: 500,
				OwnerId:          1,
			},
		}

		s.service.On("GetMachines", ctx).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := getMachineHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when invalid get machine request is made", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/machines", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := api.Message{
			Msg: "mocked error",
		}

		s.service.On("GetMachines", ctx).Return([]domain.MachineResponse{}, errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := getMachineHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})
}

func (s *HandlerTestSuite) Test_bookingHandler() {
	t := s.T()
	t.Run("when valid booking request is made", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "2021-01-01", "slots": [1,2] , "farmer_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/bookings", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := domain.NewBookingResponse{
			InvoiceId:   1,
			MachineId:   1,
			SlotsBooked: []uint{1, 2},
			TotalCost:   1000,
		}
		requestBody := domain.NewBookingRequest{
			MachineId: 1,
			Date:      "2021-01-01",
			Slots:     []uint{1, 2},
			FarmerId:  1,
		}
		s.service.On("BookMachine", ctx, requestBody).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := bookingHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when invalid booking request is made,no slots selected", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "2021-01-01", "slots": [] , "farmer_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/bookings", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))

		respBody := api.Message{
			Msg: "no slots selected",
		}

		// s.service.On("BookMachine", ctx, requestBody).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := bookingHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when invalid booking request is made,invalid slots selected", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "2021-01-01", "slots": [25, 26] , "farmer_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/bookings", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))

		respBody := api.Message{
			Msg: "invalid slot selected",
		}

		// s.service.On("BookMachine", ctx, requestBody).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := bookingHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when invalid booking request is made,invalid date format", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "13-12-2020", "slots": [1,2] , "farmer_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/bookings", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))

		respBody := api.Message{
			Msg: "invalid date",
		}

		// s.service.On("BookMachine", ctx, requestBody).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := bookingHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when error in booking ", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "2020-12-01", "slots": [1,2] , "farmer_id" : 1}`)
		r := httptest.NewRequest(http.MethodPost, "/bookings", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))

		ctx := r.Context()

		requestBody := domain.NewBookingRequest{
			MachineId: 1,
			Date:      "2020-12-01",
			Slots:     []uint{1, 2},
			FarmerId:  1,
		}
		respBody := api.Message{
			Msg: "mocked error",
		}

		s.service.On("BookMachine", ctx, requestBody).Return(domain.NewBookingResponse{}, errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := bookingHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

}

func (s *HandlerTestSuite) Test_availabilityHandler() {
	t := s.T()
	t.Run("when valid availability request is made", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "2021-01-01"}`)
		r := httptest.NewRequest(http.MethodPost, "/availability", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := domain.AvailabilityResponse{
			MachineId:      1,
			Date:           "2021-01-01",
			SlotsAvailable: []uint{1, 2},
		}
		requestBody := domain.AvailabilityRequest{
			MachineId: 1,
			Date:      "2021-01-01",
		}
		s.service.On("GetAvailability", ctx, requestBody.MachineId, requestBody.Date).Return([]uint{1, 2}, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := availabilityHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

	t.Run("when error in getting availability", func(t *testing.T) {
		bodyReader := strings.NewReader(`{"machine_id" : 1, "date" : "2021-01-01"}`)
		r := httptest.NewRequest(http.MethodPost, "/availability", (bodyReader))
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := api.Message{
			Msg: "mocked error",
		}
		requestBody := domain.AvailabilityRequest{
			MachineId: 1,
			Date:      "2021-01-01",
		}
		s.service.On("GetAvailability", ctx, requestBody.MachineId, requestBody.Date).Return([]uint{}, errors.New(
			"mocked error",
		)).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := availabilityHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

}

func (s *HandlerTestSuite) Test_getAllBookingsHandler() {
	t := s.T()
	t.Run("when valid get all booking request is made", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/bookings", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := []domain.BookingResponse{
			{
				BookingId:   1,
				MachineId:   1,
				SlotsBooked: []uint{1, 2},
			},
		}
		s.service.On("GetAllBookings", ctx, uint(1)).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := getAllBookingsHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when error in  get all booking ", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/bookings", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := api.Message{
			Msg: "mocked error",
		}
		s.service.On("GetAllBookings", ctx, uint(1)).Return([]domain.BookingResponse{}, errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := getAllBookingsHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})

}

func (s *HandlerTestSuite) Test_getAllSlotsHandler() {
	t := s.T()
	t.Run("when valid get all slots request is made", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/slots", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := []domain.SlotResponse{
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
		}

		s.service.On("GetAllSlots", ctx).Return(respBody, nil).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := getAllSlotsHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})
	t.Run("when error geting all slots request is made", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/slots", nil)
		w := httptest.NewRecorder()
		r = r.WithContext(context.WithValue(r.Context(), "token", uint(1)))
		ctx := r.Context()
		respBody := api.Message{
			Msg: "mocked error",
		}

		s.service.On("GetAllSlots", ctx).Return([]domain.SlotResponse{}, errors.New("mocked error")).Once()

		deps := dependencies{
			FarmService: s.service,
		}
		exp, _ := json.Marshal(respBody)
		got := getAllSlotsHandler(deps)
		got.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Equal(t, string(exp), w.Body.String())
	})
}
