package usecases

import (
	"github.com/leoride/tako-sim/domain"
	"math"
	"math/rand"
	"time"
)

type TripClientI interface {
	SendTripStart(*domain.Trip)
	SendTripEnd(*domain.Trip)

	SendTripData(*domain.Trip)
	SendTripComplete(*domain.Trip)

	SendRejectedAccess(*domain.DriverSwipe)
	SendDriverLate(*domain.Trip)
}

type TripService struct {
	tripClient TripClientI
	trips      []*domain.Trip
}

func NewTripService(tc TripClientI, trips []*domain.Trip) *TripService {
	ts := new(TripService)

	ts.tripClient = tc
	ts.trips = trips

	return ts
}

func (ts *TripService) HandleTripStart(t *domain.Trip) {
	t.StartTime = time.Now()
	t.OdoStart = rand.Intn(100000)
	t.Status = domain.IN_PROGRESS

	go ts.sendTripStart(t)
}

func (ts *TripService) HandleTripEnd(t *domain.Trip) {
	t.EndTime = time.Now()
	t.OdoEnd = t.OdoStart + int(math.Ceil(time.Since(t.StartTime).Hours()*float64(rand.Intn(100)+1)))
	t.Status = domain.ENDED

	go ts.sendTripEnd(t)
}

func (ts *TripService) HandleNoDrive(r *domain.Reservation) {

	t := new(domain.Trip)
	t.VehicleDevice = r.VehicleDevice
	t.AccessDevice = r.AccessDevice
	t.ReservationId = r.ReservationId
	t.Reservation = r
	t.OdoStart = 0
	t.OdoEnd = 0
	t.StartTime = time.Now()
	t.EndTime = t.StartTime
	t.Status = domain.ENDED

	r.Trip = t
	go ts.sendTripData(t)
}

func (ts *TripService) HandleTripComplete(t *domain.Trip) {
	t.Status = domain.COMPLETED

	go ts.sendTripComplete(t)
}

func (ts *TripService) HandleDriverLate(t *domain.Trip) {
	t.Status = domain.LATE

	go ts.sendDriverLate(t)
}

func (ts *TripService) HandleRejectedAccess(ds *domain.DriverSwipe) {
	go ts.sendRejectedAccess(ds)
}

func (ts *TripService) sendTripStart(t *domain.Trip) {
	time.Sleep(time.Second * 30)
	ts.tripClient.SendTripStart(t)
}

func (ts *TripService) sendTripEnd(t *domain.Trip) {
	time.Sleep(time.Second * 30)
	ts.tripClient.SendTripEnd(t)

	go ts.sendTripData(t)
}

func (ts *TripService) sendTripData(t *domain.Trip) {
	time.Sleep(time.Second * 5)
	ts.tripClient.SendTripData(t)
}

func (ts *TripService) sendTripComplete(t *domain.Trip) {
	time.Sleep(time.Second * 10)
	ts.tripClient.SendTripComplete(t)
}

func (ts *TripService) sendDriverLate(t *domain.Trip) {
	time.Sleep(time.Second * 5)
	ts.tripClient.SendDriverLate(t)
}

func (ts *TripService) sendRejectedAccess(ds *domain.DriverSwipe) {
	time.Sleep(time.Second * 30)
	ts.tripClient.SendRejectedAccess(ds)
}
