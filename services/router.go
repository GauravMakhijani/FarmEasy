package services

import (
	"net/http"

	"github.com/gorilla/mux"
)

/* The routing mechanism. Mux helps us define handler functions and the access methods */
func InitRouter(deps dependencies) (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	router.HandleFunc("/register", registerHandler(deps)).Methods(http.MethodPost)

	router.HandleFunc("/login", loginHandler(deps)).Methods(http.MethodPost)

	router.HandleFunc("/machines", ValidateUser(addMachineHandler(deps))).Methods(http.MethodPost)

	router.HandleFunc("/machines", ValidateUser(getMachineHandler(deps))).Methods(http.MethodGet)

	router.HandleFunc("/bookings", ValidateUser(bookingHandler(deps))).Methods(http.MethodPost)

	router.HandleFunc("/availability", ValidateUser(availabilityHandler(deps))).Methods(http.MethodPost)

	router.HandleFunc("/bookings/", ValidateUser(getAllBookingsHandler(deps))).Methods(http.MethodGet)

	return
}
