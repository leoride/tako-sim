package usecases

import (
	"fmt"
	"github.com/leoride/tako-sim/domain"
	"time"
)

type ReservationClientI interface {
	SendUpdate(domain.RequestI)
}

type ReservationService struct {
	reservationClient ReservationClientI
	tripService       *TripService
	reservations      []*domain.Reservation
}

type ReservationWatcherThread struct {
	TripService *TripService
	Reservation *domain.Reservation
}

func NewReservationService(rc ReservationClientI, ts *TripService, reservations []*domain.Reservation) *ReservationService {
	rs := new(ReservationService)

	rs.reservationClient = rc
	rs.tripService = ts
	rs.reservations = reservations

	return rs
}

func (rs *ReservationService) HandleNewReservation(r *domain.Reservation) {
	r.GenerateTaskNumber()
	r.TechStatus = domain.NEW

	var existingRes *domain.Reservation = nil
	for _, value := range rs.reservations {
		if value.ReservationId == r.ReservationId {
			existingRes = value
		}
	}

	if existingRes != nil {
		t := existingRes.Trip
		*existingRes = *r
		existingRes.Trip = t
		fmt.Println("Existing reservation updated:")
	} else {
		rs.reservations = append(rs.reservations, r)
		fmt.Println("New reservation received:")

		rw := new(ReservationWatcherThread)
		rw.TripService = rs.tripService
		rw.Reservation = r
		go rw.Watch()
	}

	fmt.Println(r)

	go rs.sendReservationStatusUpdates(r)
}

func (rs *ReservationService) HandleNewDriverSwipe(ds *domain.DriverSwipe) {
	ds.GenerateTaskNumber()
	ds.TechStatus = domain.NEW

	var existingRes *domain.Reservation = nil
	for _, value := range rs.reservations {

		if value.VehicleDevice.OrgaNo == ds.VehicleDevice.OrgaNo &&
			value.VehicleDevice.VehiclePhoneNo == ds.VehicleDevice.VehiclePhoneNo {

			if value.StartTime.Before(time.Now()) &&
				(time.Now().Before(value.EndTime) || value.Trip != nil && (value.Trip.Status == domain.IN_PROGRESS || value.Trip.Status == domain.LATE)) {

				if value.AccessDevice.SmartcardType == ds.AccessDevice.SmartcardType {
					switch value.AccessDevice.SmartcardType {
					case "Hitag16":
					case "Hitag32":
						if value.AccessDevice.SmartcardCardNo == ds.AccessDevice.SmartcardCardNo &&
							value.AccessDevice.SmartcardOrgaNo == ds.AccessDevice.SmartcardOrgaNo {
							existingRes = value
						}
						break
					default:
						if value.AccessDevice.SmartcardSerialNo == ds.AccessDevice.SmartcardSerialNo {
							existingRes = value
						}
						break
					}
				}
			}
		}
	}

	if existingRes == nil {
		fmt.Println("Driver swipe received, but no reservation found")
		rs.tripService.HandleRejectedAccess(ds)

	} else if existingRes != nil && existingRes.Trip == nil {
		fmt.Println("Driver swipe received, starting trip for reservation", existingRes.ReservationId)

		trip := new(domain.Trip)
		trip.VehicleDevice = existingRes.VehicleDevice
		trip.AccessDevice = existingRes.AccessDevice
		trip.ReservationId = existingRes.ReservationId
		trip.Reservation = existingRes

		existingRes.Trip = trip

		rs.tripService.HandleTripStart(trip)

	} else if existingRes != nil && existingRes.Trip != nil {
		fmt.Println("Driver swipe received for ongoing trip, ending trip", existingRes.ReservationId)
		rs.tripService.HandleTripEnd(existingRes.Trip)
	}

	go rs.sendDriverSwipeStatusUpdates(ds)
}

func (rs *ReservationService) sendReservationStatusUpdates(r *domain.Reservation) {
	time.Sleep(time.Second * 5)
	r.TechStatus = domain.SENT_TO_CUCM
	rs.reservationClient.SendUpdate(r)

	time.Sleep(time.Second * 5)
	r.TechStatus = domain.ACCEPTED_BY_CUCM
	rs.reservationClient.SendUpdate(r)

	time.Sleep(time.Second * 5)
	r.TechStatus = domain.RECEIVED
	rs.reservationClient.SendUpdate(r)
}

func (rs *ReservationService) sendDriverSwipeStatusUpdates(ds *domain.DriverSwipe) {
	time.Sleep(time.Second * 5)
	ds.TechStatus = domain.SENT_TO_CUCM
	rs.reservationClient.SendUpdate(ds)

	time.Sleep(time.Second * 5)
	ds.TechStatus = domain.ACCEPTED_BY_CUCM
	rs.reservationClient.SendUpdate(ds)

	time.Sleep(time.Second * 5)
	ds.TechStatus = domain.RECEIVED
	rs.reservationClient.SendUpdate(ds)
}

func (rw *ReservationWatcherThread) Watch() {
	for {
		r := rw.Reservation
		t := r.Trip

		if r.EndTime.Before(time.Now()) {

			if t != nil && t.Status == domain.ENDED {

				rw.TripService.HandleTripComplete(t)
				return
			} else if t != nil &&
				t.Status == domain.IN_PROGRESS &&
				r.LateAlarm == true &&
				time.Now().After(r.EndTime.Add(time.Minute*time.Duration(r.LateBuffer))) {

				rw.TripService.HandleDriverLate(t)
			} else if t == nil {

				rw.TripService.HandleNoDrive(r)
			}
		}
	}
}
