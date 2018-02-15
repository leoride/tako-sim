package usecases

import (
	"fmt"
	"github.com/google/uuid"
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

	cucmRequests map[string]*domain.DriverSwipe
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
	rs.cucmRequests = make(map[string]*domain.DriverSwipe)

	return rs
}

func (rs *ReservationService) GetReservations() []*domain.Reservation {
	return rs.reservations
}

func (rs *ReservationService) GetReservation(id string) *domain.Reservation {
	for _, value := range rs.reservations {
		if value.ReservationId == id {
			return value
		}
	}

	return nil
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
					case "Hitag_16":
					case "Hitag_32":
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

		ds.CUCMGuid = uuid.New().String()
		rs.cucmRequests[ds.CUCMGuid] = ds
		rs.tripService.HandleCUCMRequest(ds)

	} else if existingRes != nil && existingRes.Trip == nil {
		fmt.Println("Driver swipe received, starting trip for reservation", existingRes.ReservationId)

		trip := new(domain.Trip)
		trip.VehicleDevice = existingRes.VehicleDevice
		trip.AccessDevice = existingRes.AccessDevice
		trip.ReservationId = existingRes.ReservationId
		trip.Reservation = existingRes

		trip.IgnitionStatus = true
		trip.IgnitionChange = time.Now()

		existingRes.Trip = trip

		rs.tripService.HandleTripStart(trip)

	} else if existingRes != nil && existingRes.Trip != nil {
		if existingRes.Trip.Status == domain.ENDED {
			fmt.Println("Driver swipe received for ongoing trip, starting trip again", existingRes.ReservationId)

			existingRes.Trip.IgnitionStatus = true
			existingRes.Trip.IgnitionChange = time.Now()
			rs.tripService.HandleTripStart(existingRes.Trip)

		} else {
			fmt.Println("Driver swipe received for ongoing trip, ending trip", existingRes.ReservationId)
			rs.tripService.HandleTripEnd(existingRes.Trip)
		}
	}

	go rs.sendDriverSwipeStatusUpdates(ds)
}

func (rs *ReservationService) HandleNewCUCMResponse(cr *domain.CUCMResponse) {
	cr.GenerateTaskNumber()
	cr.TechStatus = domain.NEW

	ds := rs.cucmRequests[cr.Guid]

	if ds != nil {
		if cr.ReservationId == "" {
			fmt.Println("This is a refusal - Generate Rejected Access!")

			rs.tripService.HandleRejectedAccess(ds)
		} else {
			fmt.Println("This is a success - Create reservation and Generate Trip Start!")
			r := new(domain.Reservation)
			r.ReservationId = cr.ReservationId
			r.TechStatus = cr.TechStatus
			r.RequestId = cr.RequestId
			r.AccessDevice = cr.AccessDevice
			r.EndTime = cr.EndTime
			r.LateAlarm = cr.LateAlarm
			r.LateBuffer = cr.LateBuffer
			r.StartTime = cr.StartTime
			r.Timezone = cr.Timezone
			r.VehicleDevice = cr.VehicleDevice

			rs.HandleNewReservation(r)
			rs.HandleNewDriverSwipe(ds)
		}
	}

	go rs.sendCUCMResponseStatusUpdates(cr)
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

func (rs *ReservationService) sendCUCMResponseStatusUpdates(cr *domain.CUCMResponse) {
	time.Sleep(time.Second * 5)
	cr.TechStatus = domain.SENT_TO_CUCM
	rs.reservationClient.SendUpdate(cr)

	time.Sleep(time.Second * 5)
	cr.TechStatus = domain.ACCEPTED_BY_CUCM
	rs.reservationClient.SendUpdate(cr)

	time.Sleep(time.Second * 5)
	cr.TechStatus = domain.RECEIVED
	rs.reservationClient.SendUpdate(cr)
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

		if t != nil &&
			(t.Status == domain.IN_PROGRESS || t.Status == domain.LATE) &&
			t.IgnitionChange.Before(time.Now().Add(time.Duration(-5)*time.Minute)) {

			rw.TripService.HandleTripSegment(t)
		}
	}
}
