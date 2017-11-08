package domain

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	NEW              TaskStatus = "New"
	SENT_TO_CUCM     TaskStatus = "SendToCUCM"
	ACCEPTED_BY_CUCM TaskStatus = "AcceptedFromCUCM"
	RECEIVED         TaskStatus = "Done"
)

type RequestI interface {
	GetTechStatus() TaskStatus
	GetRequestId() string
	GetOrgaNo() string
	GenerateResponse() string
	GenerateStatus() string
}

type VehicleDevice struct {
	VehiclePhoneNo string `xml:"DestinationAddress>PhoneNo"`
	OrgaNo         string `xml:"OrgaNo"`
}

func generateStatus(r RequestI) string {
	return "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"\n\t<s:Body>" +
		"\n\t\t<StatusChanged xmlns=\"http://tempuri.org/\">" +
		"\n\t\t\t<status xmlns:a=\"http://invers.com\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">" +
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
		"\n\t\t\t</status>" +
		"\n\t\t</StatusChanged>" +
		"\n\t</s:Body>" +
		"\n</s:Envelope>"
}
