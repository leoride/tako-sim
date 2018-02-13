package interfaces

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/leoride/tako-sim/domain"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ReservationServiceI interface {
	HandleNewReservation(*domain.Reservation)
	HandleNewDriverSwipe(ds *domain.DriverSwipe)
	GetReservations() []*domain.Reservation
	GetReservation(id string) *domain.Reservation
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
	http.HandleFunc("/reservations/", func(w http.ResponseWriter, r *http.Request) {
		var (
			resp []byte
			err  error
		)

		id := strings.TrimPrefix(r.URL.Path, "/reservations/")

		if id == "" {
			//return all
			resp, err = json.Marshal(rl.reservationService.GetReservations())
		} else {
			//return one
			res := rl.reservationService.GetReservation(id)

			if res != nil {
				resp, err = json.Marshal(res)
			} else {
				w.WriteHeader(404)
				return
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

	http.HandleFunc("/AuthService", func(w http.ResponseWriter, r *http.Request) {

		string := "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
			"<s:Body>" +
			"<ClientLoginResponse xmlns=\"http://tempuri.org/\">" +
			"<ClientLoginResult>e691fd50-b0c2-4238-ac1e-3ac45bc75bb4</ClientLoginResult>" +
			"</ClientLoginResponse>" +
			"</s:Body>" +
			"</s:Envelope>"

		w.WriteHeader(200)
		w.Write([]byte(string))
	})

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
			} else if strings.Contains(body, "AnswerRequest") {
				resp, err = rl.listenForCUCMResponse(b)
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

func (rl *ReservationListener) listenForCUCMResponse(b []byte) ([]byte, error) {
	//TODO: do not hardcode
	response := "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"\n\t<s:Body>" +
		"\n\t\t<AnswerRequestResponse xmlns=\"http://tempuri.org/\">" +
		"\n\t\t\t<AnswerRequestResult xmlns:a=\"http://invers.com\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">" +
		"\n\t\t\t\t<a:TaskStatus>" +
		"\n\t\t\t\t\t<a:CustomerId>00000000-0000-0000-0000-000000000000</a:CustomerId>" +
		"\n\t\t\t\t\t<a:DataStatus>Sending</a:DataStatus>" +
		"\n\t\t\t\t\t<a:TaskError>NoError</a:TaskError>" +
		"\n\t\t\t\t\t<a:TaskNumber>999999</a:TaskNumber>" +
		"\n\t\t\t\t\t<a:TaskSendStatus>New</a:TaskSendStatus>" +
		"\n\t\t\t\t\t<a:Timestamp xmlns:b=\"http://schemas.datacontract.org/2004/07/Invers.DataTypes\">" +
		"\n\t\t\t\t\t\t<b:Timezone>20</b:Timezone>" +
		"\n\t\t\t\t\t\t<b:UTCDateTime>" + time.Now().UTC().Format("2006-01-02T15:04:05.0000000Z") + "</b:UTCDateTime>" +
		"\n\t\t\t\t\t</a:Timestamp>" +
		"\n\t\t\t\t\t<a:UsedCommsystem>Unknown</a:UsedCommsystem>" +
		"\n\t\t\t\t</a:TaskStatus>" +
		"\n\t\t\t</AnswerRequestResult>" +
		"\n\t\t</AnswerRequestResponse>" +
		"\n\t</s:Body>" +
		"\n</s:Envelope>"

	return []byte(response), nil
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
