package usecases

import (
	"github.com/leoride/tako-sim/domain"
	"math/rand"
	"time"
)

type TripClientI interface {
	SendTripStart(*domain.Trip)
	SendDataFobAction(*domain.Trip, bool)
	SendFirstIgnition(*domain.Trip)
	SendTripEnd(*domain.Trip)

	SendTripSegment(*domain.Trip)
	SendTripData(*domain.Trip)
	SendTripComplete(*domain.Trip)

	SendRejectedAccess(*domain.DriverSwipe)
	SendCUCMRequest(*domain.DriverSwipe)
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
	if t.OdoStart == 0 {
		t.StartTime = time.Now()
		t.OdoStart = rand.Intn(100000)
	}

	t.Status = domain.IN_PROGRESS

	go ts.sendTripStart(t)
}

func (ts *TripService) HandleTripEnd(t *domain.Trip) {
	t.EndTime = time.Now()
	//t.OdoEnd = t.OdoStart + int(math.Ceil(time.Since(t.StartTime).Hours()*float64(rand.Intn(100)+1)))
	t.Status = domain.ENDED

	if t.IgnitionStatus == true {
		ts.HandleTripSegment(t)
	}

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

func (ts *TripService) HandleTripSegment(t *domain.Trip) {
	t.IgnitionStatus = !t.IgnitionStatus
	t.IgnitionChange = time.Now()

	if t.IgnitionStatus == false {
		if t.OdoEnd == 0 {
			t.OdoEnd = t.OdoStart
		}
		t.OdoEnd = t.OdoEnd + 3
	}
	t.EndTime = time.Now()

	go ts.sendTripSegment(t)
}

func (ts *TripService) HandleDriverLate(t *domain.Trip) {
	t.Status = domain.LATE

	go ts.sendDriverLate(t)
}

func (ts *TripService) HandleRejectedAccess(ds *domain.DriverSwipe) {
	go ts.sendRejectedAccess(ds)
}

func (ts *TripService) HandleCUCMRequest(ds *domain.DriverSwipe) {
	go ts.sendCUCMRequest(ds)
}

func (ts *TripService) sendTripStart(t *domain.Trip) {
	time.Sleep(time.Second * 30)
	ts.tripClient.SendTripStart(t)

	time.Sleep(time.Second * 10)
	ts.tripClient.SendDataFobAction(t, true)

	if t.OdoEnd == 0 {
		time.Sleep(time.Second * 10)
		ts.tripClient.SendFirstIgnition(t)
	}
}

func (ts *TripService) sendTripEnd(t *domain.Trip) {
	time.Sleep(time.Second * 30)
	ts.tripClient.SendTripEnd(t)

	time.Sleep(time.Second * 10)
	ts.tripClient.SendDataFobAction(t, false)

	go ts.sendTripData(t)
}

func (ts *TripService) sendTripSegment(t *domain.Trip) {
	time.Sleep(time.Second * 5)
	ts.tripClient.SendTripSegment(t)
}

func (ts *TripService) sendTripData(t *domain.Trip) {
	time.Sleep(time.Second * 5)
	ts.tripClient.SendTripData(t)
}

func (ts *TripService) sendTripComplete(t *domain.Trip) {
	time.Sleep(time.Second * 60)
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

func (ts *TripService) sendCUCMRequest(ds *domain.DriverSwipe) {
	time.Sleep(time.Second * 30)
	ts.tripClient.SendCUCMRequest(ds)
}
