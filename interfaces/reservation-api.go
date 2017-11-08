package interfaces

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/leoride/tako-sim/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

type ReservationServiceI interface {
	HandleNewReservation(*domain.Reservation)
	HandleNewDriverSwipe(ds *domain.DriverSwipe)
}

type ReservationListener struct {
	reservationService ReservationServiceI
}

type ReservationClient struct {
	takoEndpoint string
}

func NewReservationClient(takoEndpoint string) *ReservationClient {
	rc := new(ReservationClient)

	rc.takoEndpoint = takoEndpoint

	return rc
}

func NewReservationListener(rs ReservationServiceI) *ReservationListener {
	rl := new(ReservationListener)
	rl.reservationService = rs

	return rl
}

func (rl *ReservationListener) Listen() {
	http.HandleFunc("/ComService", func(w http.ResponseWriter, r *http.Request) {
		var (
			b    []byte
			resp []byte
			err  error
		)

		if b, err = ioutil.ReadAll(r.Body); err == nil {
			body := string(b)

			if strings.Contains(body, "SendReservation") {
				resp, err = rl.listenForReservation(b)
			} else if strings.Contains(body, "SendVirtualSmartCard") {
				resp, err = rl.listenForSwipe(b)
			} else {
				err = fmt.Errorf("Unsupported method")
			}
		}

		if err != nil {
			fmt.Println("ERROR:", err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(200)
			w.Write(resp)
		}
	})
}

func (rl *ReservationListener) listenForReservation(b []byte) ([]byte, error) {
	rt := new(domain.Reservation)

	if err := xml.Unmarshal(b, rt); err == nil {
		rl.reservationService.HandleNewReservation(rt)
		response := rt.GenerateResponse()

		return []byte(response), nil

	} else {
		return nil, fmt.Errorf("Error processing request:", err)
	}
}

func (rl *ReservationListener) listenForSwipe(b []byte) ([]byte, error) {
	ds := new(domain.DriverSwipe)

	if err := xml.Unmarshal(b, ds); err == nil {
		rl.reservationService.HandleNewDriverSwipe(ds)
		response := ds.GenerateResponse()

		return []byte(response), nil

	} else {
		return nil, fmt.Errorf("Error processing request:", err)
	}
}

func (rc *ReservationClient) SendUpdate(r domain.RequestI) {
	body := []byte(r.GenerateStatus())
	req, err := http.NewRequest("POST", rc.takoEndpoint+"/ws/invers/21/"+r.GetOrgaNo()+"/com", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("status update sent:", r.GetTechStatus())
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Status update error:", err)
	}
}
