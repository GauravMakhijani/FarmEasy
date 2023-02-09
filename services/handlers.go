package services

import (
	"FarmEasy/api"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MsgResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func registerHandler(deps dependencies) http.HandlerFunc {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var farmer NewFarmer
		err := json.NewDecoder(req.Body).Decode(&farmer)
		if err != nil {
			api.Response(rw, http.StatusBadRequest, err.Error())
			return
		}

		if err = ValidateFarmerPhone(farmer.Phone); err != nil {
			api.Response(rw, http.StatusBadRequest, err.Error())

			return
		}
		if err = ValidateFarmerEmail(farmer.Email); err != nil {

			api.Response(rw, http.StatusBadRequest, err.Error())

			return
		}

		addedFarmer, err := deps.FarmService.Register(req.Context(), farmer)

		if err != nil {
			api.Response(rw, http.StatusBadRequest, err.Error())

			return
		}
		api.Response(rw, http.StatusCreated, addedFarmer)

	})
}

func loginHandler(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fAuth NewLogin

		if err := json.NewDecoder(r.Body).Decode(&fAuth); err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ValidateFarmerEmail(fAuth.Email); err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())

			return
		}
		tokenString, err := deps.FarmService.Login(r.Context(), fAuth)

		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())

			return
		}

		msg := MsgResponse{Message: "Login Successful", Token: tokenString}

		api.Response(w, http.StatusOK, msg)

	}
}

func addMachineHandler(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var machine NewMachine
		farmerId := r.Context().Value("token")

		machine.OwnerId = farmerId.(uint)
		if err := json.NewDecoder(r.Body).Decode(&machine); err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}
		addedMachine, err := deps.FarmService.AddMachine(r.Context(), machine)

		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())

			return
		}

		api.Response(w, http.StatusOK, addedMachine)
	}
}

func getMachineHandler(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		machines, err := deps.FarmService.GetMachines(r.Context())

		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}

		api.Response(w, http.StatusOK, machines)

	}
}

func bookingHandler(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var booking NewBooking
		id := r.Context().Value("token")
		farmerId := id.(uint)
		booking.FarmerId = farmerId
		if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := ValidateBookingslots(booking.Slots); err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())

			return
		}
		if err := ValidateBookingDate(booking.Date); err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}

		addedBooking, err := deps.FarmService.BookMachine(r.Context(), booking)

		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}
		api.Response(w, http.StatusCreated, addedBooking)
	}
}

func availabilityHandler(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		machineId := vars["id"]
		id, err := strconv.Atoi(machineId)
		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}
		slotsAvailable, err := deps.FarmService.GetAvailability(r.Context(), uint(id))
		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}
		api.Response(w, http.StatusOK, slotsAvailable)
	}
}

func getAllBookingsHandler(deps dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("token")
		farmerId := id.(uint)
		bookings, err := deps.FarmService.GetAllBookings(r.Context(), uint(farmerId))
		if err != nil {
			api.Response(w, http.StatusBadRequest, err.Error())
			return
		}
		api.Response(w, http.StatusOK, bookings)
	}
}
