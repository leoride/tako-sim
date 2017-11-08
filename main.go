package main

import (
	"flag"
	"fmt"
	"github.com/leoride/tako-sim/domain"
	"github.com/leoride/tako-sim/interfaces"
	"github.com/leoride/tako-sim/usecases"
	"log"
	"net/http"
)

func main() {
	fmt.Println("===== STARTING TAKO TECH SIMULATOR =====")

	var (
		takoEndpoint string
		port         int

		tc *interfaces.TripClient
		ts *usecases.TripService

		rc *interfaces.ReservationClient
		rs *usecases.ReservationService
		rl *interfaces.ReservationListener

		reservations []*domain.Reservation = make([]*domain.Reservation, 0)
		trips        []*domain.Trip        = make([]*domain.Trip, 0)
	)

	takoEndpoint = *flag.String("takoEndpoint", "http://localhost:8080/tako-fc", "Tako FC root URL")
	port = *flag.Int("port", 8282, "Port the app listens to")

	tc = interfaces.NewTripClient(takoEndpoint)
	ts = usecases.NewTripService(tc, trips)

	rc = interfaces.NewReservationClient(takoEndpoint)
	rs = usecases.NewReservationService(rc, ts, reservations)
	rl = interfaces.NewReservationListener(rs)

	rl.Listen()

	log.Panic(http.ListenAndServe(":"+fmt.Sprint(port), nil))
}
