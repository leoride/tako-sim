package interfaces

import (
	"bytes"
	"fmt"
	"github.com/leoride/tako-sim/domain"
	"net/http"
)

type TripClient struct {
	takoEndpoint string
}

func NewTripClient(takoEndpoint string) *TripClient {
	tc := new(TripClient)

	tc.takoEndpoint = takoEndpoint

	return tc
}

func (tc *TripClient) SendTripStart(t *domain.Trip) {
	body := []byte(t.GenerateTripStart())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("trip start sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Trip start error:", err)
	}
}

func (tc *TripClient) SendFirstIgnition(t *domain.Trip) {
	body := []byte(t.GenerateFirstIgnition())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("first ignition sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("first ignition error:", err)
	}
}

func (tc *TripClient) SendDataFobAction(t *domain.Trip, removed bool) {
	body := []byte(t.GenerateDataFobAction(removed))
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("datafob event sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("datafob event error:", err)
	}
}

func (tc *TripClient) SendTripEnd(t *domain.Trip) {
	body := []byte(t.GenerateTripEnd())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("trip end sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Trip end error:", err)
	}
}

func (tc *TripClient) SendTripSegment(t *domain.Trip) {
	body := []byte(t.GenerateTripSegment())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/trip", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("trip segment sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Trip segment error:", err)
	}
}

func (tc *TripClient) SendTripData(t *domain.Trip) {
	body := []byte(t.GenerateTripData())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/trip", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("trip data sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Trip data error:", err)
	}
}

func (tc *TripClient) SendTripComplete(t *domain.Trip) {
	body := []byte(t.GenerateTripComplete())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("trip complete sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Trip complete error:", err)
	}
}

func (tc *TripClient) SendDriverLate(t *domain.Trip) {
	body := []byte(t.GenerateDriverLate())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+t.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("driver late sent:", fmt.Sprint(t.Status))
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Driver late error:", err)
	}
}

func (tc *TripClient) SendRejectedAccess(ds *domain.DriverSwipe) {
	body := []byte(ds.GenerateRejectedAccess())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+ds.VehicleDevice.OrgaNo+"/event", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("rejected access sent")
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("Rejected access error:", err)
	}
}

func (tc *TripClient) SendCUCMRequest(ds *domain.DriverSwipe) {
	body := []byte(ds.GenerateCUCMRequest())
	req, err := http.NewRequest("POST", tc.takoEndpoint+"/ws/invers/21/"+ds.VehicleDevice.OrgaNo+"/res", bytes.NewBuffer(body))

	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)

		if err == nil {
			fmt.Println("CUCM request sent")
			fmt.Println("response Status:", resp.Status)
			defer resp.Body.Close()
		}
	}

	if err != nil {
		fmt.Println("CUCM request error:", err)
	}
}
