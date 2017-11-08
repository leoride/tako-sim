package domain

import (
	"fmt"
	"math/rand"
	"time"
)

type Reservation struct {
	Timezone      int `xml:"Body>SendReservation>task>Reservation>Start>Timezone"`
	TechStatus    TaskStatus
	VehicleDevice VehicleDevice `xml:"Body>SendReservation>task>Destination"`
	AccessDevice  AccessDevice  `xml:"Body>SendReservation>task>Reservation>UserAccessList>UserAccess"`
	ReservationId string        `xml:"Body>SendReservation>task>Reservation>ReservationNo"`
	RequestId     string        `xml:"Body>SendReservation>task>TaskNumber"`
	StartTime     time.Time     `xml:"Body>SendReservation>task>Reservation>Start>UTCDateTime"`
	EndTime       time.Time     `xml:"Body>SendReservation>task>Reservation>Stop>UTCDateTime"`
	LateAlarm     bool          `xml:"Body>SendReservation>task>Reservation>ReturnOptions>DelayMessage"`
	LateBuffer    int           `xml:"Body>SendReservation>task>Reservation>ReturnOptions>DelayTime"`
	Trip          *Trip
}

type AccessDevice struct {
	SmartcardSerialNo string `xml:"SerialNo"`
	SmartcardCardNo   string `xml:"CardNo"`
	SmartcardOrgaNo   string `xml:"CardOrga"`
	SmartcardType     string `xml:"Type"`
}

func (r *Reservation) GetTimezone() *time.Location {
	var loc *time.Location

	switch r.Timezone {
	case 2:
		loc, _ = time.LoadLocation("Pacific/Honolulu")
		break
	case 4:
		loc, _ = time.LoadLocation("PST8PDT")
		break
	case 10:
		loc, _ = time.LoadLocation("MST7MDT")
		break
	case 20:
		loc, _ = time.LoadLocation("CST6CDT")
		break
	case 35:
		loc, _ = time.LoadLocation("EST5EDT")
		break
	case 85:
		loc, _ = time.LoadLocation("Europe/London")
		break
	}

	return loc
}

func (r *Reservation) GetTechStatus() TaskStatus {
	return r.TechStatus
}

func (r *Reservation) GetRequestId() string {
	return r.RequestId
}

func (r *Reservation) GetOrgaNo() string {
	return r.VehicleDevice.OrgaNo
}

func (r *Reservation) GenerateStatus() string {
	return generateStatus(RequestI(r))
}

func (r *Reservation) GenerateTaskNumber() {
	r.RequestId = fmt.Sprint(rand.Intn(1000000))
}

func (rt *Reservation) String() string {
	string := ""

	string += fmt.Sprint("Reservation Request received:") +
		fmt.Sprint("\n\t") + fmt.Sprint("- SmartcardSerialNo: ") + fmt.Sprintf("%s", rt.AccessDevice.SmartcardSerialNo) +
		fmt.Sprint("\n\t") + fmt.Sprint("- ReservationId: ") + fmt.Sprintf("%s", rt.ReservationId) +
		fmt.Sprint("\n\t") + fmt.Sprint("- StartTime: ") + fmt.Sprintf("%s", rt.StartTime) +
		fmt.Sprint("\n\t") + fmt.Sprint("- EndTime: ") + fmt.Sprintf("%s", rt.EndTime) +
		fmt.Sprint("\n\t") + fmt.Sprint("- VehiclePhoneNo: ") + fmt.Sprintf("%s", rt.VehicleDevice.VehiclePhoneNo) +
		fmt.Sprint("\n\t") + fmt.Sprint("- OrgaNo: ") + fmt.Sprintf("%s", rt.VehicleDevice.OrgaNo) +
		fmt.Sprint("\n\t") + fmt.Sprint("- RequestId: ") + fmt.Sprintf("%s", rt.RequestId)

	return string
}

func (r *Reservation) GenerateResponse() string {
	return "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"\n\t<s:Body>" +
		"\n\t\t<SendReservationResponse xmlns=\"http://tempuri.org/\">" +
		"\n\t\t\t<SendReservationResult xmlns:a=\"http://invers.com\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">" +
		"\n\t\t\t\t<a:CustomerId>00000000-0000-0000-0000-000000000000</a:CustomerId>" +
		"\n\t\t\t\t<a:DataStatus>Sending</a:DataStatus>" +
		"\n\t\t\t\t<a:TaskError>NoError</a:TaskError>" +
		"\n\t\t\t\t<a:TaskNumber>" + r.GetRequestId() + "</a:TaskNumber>" +
		"\n\t\t\t\t<a:TaskSendStatus>" + fmt.Sprint(r.GetTechStatus()) + "</a:TaskSendStatus>" +
		"\n\t\t\t\t<a:Timestamp xmlns:b=\"http://schemas.datacontract.org/2004/07/Invers.DataTypes\">" +
		"\n\t\t\t\t\t<b:Timezone>20</b:Timezone>" +
		"\n\t\t\t\t\t<b:UTCDateTime>" + time.Now().UTC().Format("2006-01-02T15:04:05.0000000Z") + "</b:UTCDateTime>" +
		"\n\t\t\t\t</a:Timestamp>" +
		"\n\t\t\t\t<a:UsedCommsystem>Unknown</a:UsedCommsystem>" +
		"\n\t\t\t</SendReservationResult>" +
		"\n\t\t</SendReservationResponse>" +
		"\n\t</s:Body>" +
		"\n</s:Envelope>"
}
