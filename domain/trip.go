package domain

import (
	"fmt"
	"math/rand"
	"time"
)

type TripStatus string
type EventName string

const (
	IN_PROGRESS TripStatus = "STARTED"
	LATE        TripStatus = "LATE"
	ENDED       TripStatus = "ENDED"
	COMPLETED   TripStatus = "COMPLETED"

	TRIP_START    EventName = "TripStartFromDevice"
	TRIP_END      EventName = "TripEndFromDevice"
	TRIP_COMPLETE EventName = "TripFinished"

	REJECTED_ACCESS EventName = "RejectedAccess"
	LATE_DRIVER     EventName = "DelayedTripEnd"
)

type Trip struct {
	Reservation   *Reservation `json:"-"`
	VehicleDevice VehicleDevice
	AccessDevice  AccessDevice
	ReservationId string
	TripId        string
	StartTime     time.Time
	EndTime       time.Time
	OdoStart      int
	OdoEnd        int
	Status        TripStatus
}

type DriverSwipe struct {
	TechStatus    TaskStatus
	RequestId     string              `xml:"Body>SendVirtualSmartCard>task>TaskNumber"`
	VehicleDevice VehicleDevice       `xml:"Body>SendVirtualSmartCard>task>Destination"`
	AccessDevice  VirtualAccessDevice `xml:"Body>SendVirtualSmartCard>task>VirtualSmartCard"`
}

type VirtualAccessDevice struct {
	SmartcardSerialNo string `xml:"CocosNumber"`
	SmartcardCardNo   string `xml:"UserNumber"`
	SmartcardOrgaNo   string `xml:"OrgaRef"`
	SmartcardType     string `xml:"Type"`
}

func (r *DriverSwipe) GetTechStatus() TaskStatus {
	return r.TechStatus
}

func (r *DriverSwipe) GetRequestId() string {
	return r.RequestId
}

func (r *DriverSwipe) GetOrgaNo() string {
	return r.VehicleDevice.OrgaNo
}

func (r *DriverSwipe) GenerateStatus() string {
	return generateStatus(RequestI(r))
}

func (r *DriverSwipe) GenerateTaskNumber() {
	r.RequestId = fmt.Sprint(rand.Intn(1000000))
}

func (ds *DriverSwipe) GenerateResponse() string {
	return "<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"\n\t<s:Body>" +
		"\n\t\t<SendVirtualSmartCardResponse xmlns=\"http://tempuri.org/\">" +
		"\n\t\t\t<SendVirtualSmartCardResult xmlns:a=\"http://invers.com\" xmlns:i=\"http://www.w3.org/2001/XMLSchema-instance\">" +
		"\n\t\t\t\t<a:CustomerId>00000000-0000-0000-0000-000000000000</a:CustomerId>" +
		"\n\t\t\t\t<a:DataStatus>Sending</a:DataStatus>" +
		"\n\t\t\t\t<a:TaskError>NoError</a:TaskError>" +
		"\n\t\t\t\t<a:TaskNumber>" + ds.GetRequestId() + "</a:TaskNumber>" +
		"\n\t\t\t\t<a:TaskSendStatus>" + fmt.Sprint(ds.GetTechStatus()) + "</a:TaskSendStatus>" +
		"\n\t\t\t\t<a:Timestamp xmlns:b=\"http://schemas.datacontract.org/2004/07/Invers.DataTypes\">" +
		"\n\t\t\t\t\t<b:Timezone>20</b:Timezone>" +
		"\n\t\t\t\t\t<b:UTCDateTime>" + time.Now().UTC().Format("2006-01-02T15:04:05.0000000Z") + "</b:UTCDateTime>" +
		"\n\t\t\t\t</a:Timestamp>" +
		"\n\t\t\t\t<a:UsedCommsystem>Unknown</a:UsedCommsystem>" +
		"\n\t\t\t</SendVirtualSmartCardResult>" +
		"\n\t\t</SendVirtualSmartCardResponse>" +
		"\n\t</s:Body>" +
		"\n</s:Envelope>"
}

func (ds *DriverSwipe) GenerateRejectedAccess() string {
	return generateProblemEvent(REJECTED_ACCESS, nil, ds)
}

func (t *Trip) GenerateDriverLate() string {
	return generateProblemEvent(LATE_DRIVER, t, nil)
}

func (t *Trip) GenerateTripStart() string {
	return t.generateEvent(TRIP_START)
}

func (t *Trip) GenerateTripEnd() string {
	return t.generateEvent(TRIP_END)
}

func (t *Trip) GenerateTripData() string {
	var tripData string

	didNotDrive := "false"
	if t.StartTime.Equal(t.EndTime) && t.OdoStart == t.OdoEnd {
		didNotDrive = "true"
	}

	loc := t.Reservation.GetTimezone()

	tripData = "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"	<soap:Body>" +
		"		<ns4:RawTripEvaluated xmlns=\"http://schemas.datacontract.org/2004/07/Invers.Ics.Interface\" xmlns:ns2=\"http://invers.com\" xmlns:ns3=\"http://schemas.datacontract.org/2004/07/Invers.DataTypes\" xmlns:ns4=\"http://tempuri.org/\" xmlns:ns5=\"http://schemas.microsoft.com/2003/10/Serialization/\" xmlns:ns6=\"http://schemas.datacontract.org/2004/07/System.Net.Mail\">" +
		"			<ns4:trip>" +
		"				<ns2:AdditionalParameters>" +
		"					<list>" +
		"						<AdditionalParameter>" +
		"							<Name>KeyStatus</Name>" +
		"							<ParamType>Int32</ParamType>" +
		"							<Value>4</Value>" +
		"						</AdditionalParameter>" +
		"					</list>" +
		"				</ns2:AdditionalParameters>" +
		"				<ns2:AdjustmentDistance>0</ns2:AdjustmentDistance>" +
		"				<ns2:Complete>true</ns2:Complete>" +
		"				<ns2:ComputedDrivingDistance>0.0</ns2:ComputedDrivingDistance>" +
		"				<ns2:ComputedStartMileage>" + fmt.Sprint(t.OdoStart) + "</ns2:ComputedStartMileage>" +
		"				<ns2:ComputedStopMileage>" + fmt.Sprint(t.OdoEnd) + "</ns2:ComputedStopMileage>" +
		"				<ns2:DistanceConversionFactor>1.0</ns2:DistanceConversionFactor>" +
		"				<ns2:DrivingDistance>0</ns2:DrivingDistance>" +
		"				<ns2:EmergencyReason>NoEmergencyTrip</ns2:EmergencyReason>" +
		"				<ns2:EmergencyTrip>false</ns2:EmergencyTrip>" +
		"				<ns2:Fuel>100</ns2:Fuel>" +
		"				<ns2:Illegal>false</ns2:Illegal>" +
		"				<ns2:NewTrip>true</ns2:NewTrip>" +
		"				<ns2:ReservationItem>" +
		"					<ID>0</ID>" +
		"					<Name/>" +
		"					<OrgaNo>" + t.VehicleDevice.OrgaNo + "</OrgaNo>" +
		"					<Type>BCSA</Type>" +
		"				</ns2:ReservationItem>" +
		"				<ns2:ReservationNo>" + t.ReservationId + "</ns2:ReservationNo>" +
		"				<ns2:ReservationType>0</ns2:ReservationType>" +
		"				<ns2:SentStatus>Sending</ns2:SentStatus>" +
		"				<ns2:Source>" +
		"					<ns2:DestinationAddress>" +
		"						<ns2:Fax/>" +
		"						<ns2:MailAddress>" +
		"							<ns2:BCC/>" +
		"							<ns2:CC/>" +
		"							<ns2:From>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns2:From>" +
		"							<ns2:Password/>" +
		"							<ns2:Priority>Normal</ns2:Priority>" +
		"							<ns2:ReplyTo>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns2:ReplyTo>" +
		"							<ns2:Sender>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns2:Sender>" +
		"							<ns2:Server/>" +
		"							<ns2:To/>" +
		"							<ns2:UserName/>" +
		"						</ns2:MailAddress>" +
		"						<ns2:PhoneNo>+" + t.VehicleDevice.VehiclePhoneNo + "</ns2:PhoneNo>" +
		"						<ns2:SIM>" +
		"							<GSMDataNo/>" +
		"							<GSMFaxNo/>" +
		"							<GSMProvider/>" +
		"							<GSMVoiceNo/>" +
		"							<ID/>" +
		"						</ns2:SIM>" +
		"						<ns2:TCPHost/>" +
		"						<ns2:TCPPort>0</ns2:TCPPort>" +
		"					</ns2:DestinationAddress>" +
		"					<ns2:DestinationType>IBOXX</ns2:DestinationType>" +
		"					<ns2:Firmwareversion/>" +
		"					<ns2:OrgaNo>" + t.VehicleDevice.OrgaNo + "</ns2:OrgaNo>" +
		"					<ns2:SourceNo>95539211389632515</ns2:SourceNo>" +
		"				</ns2:Source>" +
		"				<ns2:Start>" + t.StartTime.In(loc).Format("2006-01-02T15:04:05") + "</ns2:Start>" +
		"				<ns2:StartGPS>" +
		"					<ns2:Altitude>0.0</ns2:Altitude>" +
		"					<ns2:Distance>0</ns2:Distance>" +
		"					<ns2:Format>ddd_dddddd</ns2:Format>" +
		"					<ns2:Latitude>51.493905</ns2:Latitude>" +
		"					<ns2:LatitudeHemisphere>32</ns2:LatitudeHemisphere>" +
		"					<ns2:Longitude>-0.10749166666666667</ns2:Longitude>" +
		"					<ns2:LongitudeHemisphere>32</ns2:LongitudeHemisphere>" +
		"					<ns2:Quality>1</ns2:Quality>" +
		"					<ns2:SatInUse>8</ns2:SatInUse>" +
		"					<ns2:Timestamp>2015-04-28T04:49:43</ns2:Timestamp>" +
		"				</ns2:StartGPS>" +
		"				<ns2:StartMileage>" + fmt.Sprint(t.OdoStart) + "</ns2:StartMileage>" +
		"				<ns2:Stop>" + t.EndTime.In(loc).Format("2006-01-02T15:04:05") + "</ns2:Stop>" +
		"				<ns2:StopGPS>" +
		"					<ns2:Altitude>0.0</ns2:Altitude>" +
		"					<ns2:Distance>0</ns2:Distance>" +
		"					<ns2:Format>ddd_dddddd</ns2:Format>" +
		"					<ns2:Latitude>51.49389166666666</ns2:Latitude>" +
		"					<ns2:LatitudeHemisphere>32</ns2:LatitudeHemisphere>" +
		"					<ns2:Longitude>-0.107495</ns2:Longitude>" +
		"					<ns2:LongitudeHemisphere>32</ns2:LongitudeHemisphere>" +
		"					<ns2:Quality>1</ns2:Quality>" +
		"					<ns2:SatInUse>9</ns2:SatInUse>" +
		"					<ns2:Timestamp>2015-05-01T06:51:54</ns2:Timestamp>" +
		"				</ns2:StopGPS>" +
		"				<ns2:StopMileage>" + fmt.Sprint(t.OdoEnd) + "</ns2:StopMileage>" +
		"				<ns2:SystemTimestamp>" +
		"					<ns3:Timezone>20</ns3:Timezone>" +
		"					<ns3:UTCDateTime>" + time.Now().UTC().Format("2006-01-02T15:04:05.0000000Z") + "</ns3:UTCDateTime>" + //2015-05-01T07:20:20.2299095-05:00
		"				</ns2:SystemTimestamp>" +
		"				<ns2:Tlv/>" +
		"				<ns2:TripNo>1</ns2:TripNo>" +
		"				<ns2:Unused>" + didNotDrive + "</ns2:Unused>" +
		"				<ns2:UserAccess>" +
		"					<ns2:CardExtension>32</ns2:CardExtension>" +
		"					<ns2:CardNo>" + t.AccessDevice.SmartcardCardNo + "</ns2:CardNo>" +
		"					<ns2:CardOrga>" + t.AccessDevice.SmartcardOrgaNo + "</ns2:CardOrga>" +
		"					<ns2:CocosSerialNo>0</ns2:CocosSerialNo>" +
		"					<ns2:PIN/>" +
		"					<ns2:PINs>0</ns2:PINs>" +
		"					<ns2:SerialNo>" + t.AccessDevice.SmartcardSerialNo + "</ns2:SerialNo>" +
		"					<ns2:TAN>0</ns2:TAN>" +
		"					<ns2:TempPIN/>" +
		"					<ns2:Type>" + t.AccessDevice.SmartcardType + "</ns2:Type>" +
		"				</ns2:UserAccess>" +
		"				<ns2:WithoutReservation>false</ns2:WithoutReservation>" +
		"			</ns4:trip>" +
		"		</ns4:RawTripEvaluated>" +
		"	</soap:Body>" +
		"</soap:Envelope>"

	return tripData
}

func (t *Trip) GenerateTripComplete() string {
	return t.generateEvent(TRIP_COMPLETE)
}

func (t *Trip) generateEvent(en EventName) string {
	loc := t.Reservation.GetTimezone()
	var mil string

	if en == TRIP_START {
		mil = fmt.Sprint(t.OdoStart)
	} else {
		mil = fmt.Sprint(t.OdoEnd)
	}

	return "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"	<soap:Body>" +
		"		<ns5:UsageEventReceived xmlns=\"http://schemas.datacontract.org/2004/07/Invers.Ics.Interface\" xmlns:ns2=\"http://schemas.datacontract.org/2004/07/Invers.Ics.Interface.EvMo\" xmlns:ns3=\"http://invers.com\" xmlns:ns4=\"http://schemas.datacontract.org/2004/07/Invers.DataTypes\" xmlns:ns5=\"http://tempuri.org/\" xmlns:ns6=\"http://schemas.microsoft.com/2003/10/Serialization/\" xmlns:ns7=\"http://schemas.datacontract.org/2004/07/System.Net.Mail\">" +
		"			<ns5:usage>" +
		"				<ns2:AdditionalParameters>" +
		"					<list>" +
		"						<AdditionalParameter>" +
		"							<Name>ReservationNo</Name>" +
		"							<ParamType>UInt32</ParamType>" +
		"							<Value>" + t.ReservationId + "</Value>" +
		"						</AdditionalParameter>" +
		"					</list>" +
		"				</ns2:AdditionalParameters>" +
		"				<ns2:Description>" + fmt.Sprint(en) + "</ns2:Description>" +
		"				<ns2:Id>27642813</ns2:Id>" +
		"				<ns2:Position>" +
		"					<ns3:Altitude>0.0</ns3:Altitude>" +
		"					<ns3:Distance>0</ns3:Distance>" +
		"					<ns3:Format>ddd_dddddd</ns3:Format>" +
		"					<ns3:Latitude>51.49765166666667</ns3:Latitude>" +
		"					<ns3:LatitudeHemisphere>32</ns3:LatitudeHemisphere>" +
		"					<ns3:Longitude>-0.217075</ns3:Longitude>" +
		"					<ns3:LongitudeHemisphere>32</ns3:LongitudeHemisphere>" +
		"					<ns3:Quality>1</ns3:Quality>" +
		"					<ns3:SatInUse>8</ns3:SatInUse>" +
		"					<ns3:Timestamp>" + time.Now().In(loc).Format("2006-01-02T15:04:05") + "</ns3:Timestamp>" +
		"				</ns2:Position>" +
		"				<ns2:SentStatus>Sending</ns2:SentStatus>" +
		"				<ns2:Source>" +
		"					<ns3:DestinationAddress>" +
		"						<ns3:Fax/>" +
		"						<ns3:MailAddress>" +
		"							<ns3:BCC/>" +
		"							<ns3:CC/>" +
		"							<ns3:From>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns3:From>" +
		"							<ns3:Password/>" +
		"							<ns3:Priority>Normal</ns3:Priority>" +
		"							<ns3:ReplyTo>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns3:ReplyTo>" +
		"							<ns3:Sender>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns3:Sender>" +
		"							<ns3:Server/>" +
		"							<ns3:To/>" +
		"							<ns3:UserName/>" +
		"						</ns3:MailAddress>" +
		"						<ns3:PhoneNo>+" + t.VehicleDevice.VehiclePhoneNo + "</ns3:PhoneNo>" +
		"						<ns3:SIM>" +
		"							<GSMDataNo/>" +
		"							<GSMFaxNo/>" +
		"							<GSMProvider/>" +
		"							<GSMVoiceNo/>" +
		"							<ID/>" +
		"						</ns3:SIM>" +
		"						<ns3:TCPHost/>" +
		"						<ns3:TCPPort>0</ns3:TCPPort>" +
		"					</ns3:DestinationAddress>" +
		"					<ns3:DestinationType>BCSA</ns3:DestinationType>" +
		"					<ns3:Firmwareversion/>" +
		"					<ns3:OrgaNo>" + t.VehicleDevice.OrgaNo + "</ns3:OrgaNo>" +
		"					<ns3:SourceNo>132309508675338243</ns3:SourceNo>" +
		"				</ns2:Source>" +
		"				<ns2:SystemTimestamp xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:nil=\"true\"/>" +
		"				<ns2:Timestamp>" + time.Now().In(loc).Format("2006-01-02T15:04:05") + "</ns2:Timestamp>" +
		"				<ns2:Tlv/>" +
		"				<ns2:Type>12</ns2:Type>" +
		"				<ns3:AnswerList/>" +
		"				<ns3:BcStatus>WaitingForPIN</ns3:BcStatus>" +
		"				<ns3:CallReason>Unknown</ns3:CallReason>" +
		"				<ns3:CentralLockState>" +
		"					<NewOpen>false</NewOpen>" +
		"					<OpenCmd>false</OpenCmd>" +
		"					<Reason>Card</Reason>" +
		"				</ns3:CentralLockState>" +
		"				<ns3:DataFob>0</ns3:DataFob>" +
		"				<ns3:Driver>false</ns3:Driver>" +
		"				<ns3:DrivingDistance>0</ns3:DrivingDistance>" +
		"				<ns3:EnterPassengerCount>0</ns3:EnterPassengerCount>" +
		"				<ns3:Fuel>-1</ns3:Fuel>" +
		"				<ns3:FuelCard>0</ns3:FuelCard>" +
		"				<ns3:LedStatus>" +
		"					<Green>false</Green>" +
		"					<Red>false</Red>" +
		"					<Yellow>false</Yellow>" +
		"				</ns3:LedStatus>" +
		"				<ns3:Mileage>" + mil + "</ns3:Mileage>" +
		"				<ns3:PassengerCount>0</ns3:PassengerCount>" +
		"				<ns3:Pause>false</ns3:Pause>" +
		"				<ns3:PinData>" +
		"					<PINs>1</PINs>" +
		"					<Result>OK</Result>" +
		"					<Tries>1</Tries>" +
		"				</ns3:PinData>" +
		"				<ns3:ReservationItem>" +
		"					<ID>0</ID>" +
		"					<Name/>" +
		"					<OrgaNo>" + t.VehicleDevice.OrgaNo + "</OrgaNo>" +
		"					<Type>BCSA</Type>" +
		"				</ns3:ReservationItem>" +
		"				<ns3:ReservationTypeId>-1</ns3:ReservationTypeId>" +
		"				<ns3:SpeedAlert>" +
		"					<Delay>PT0S</Delay>" +
		"					<Limit>0</Limit>" +
		"					<Speed>0</Speed>" +
		"				</ns3:SpeedAlert>" +
		"				<ns3:Start>1900-01-01T00:00:00</ns3:Start>" +
		"				<ns3:Stop>1900-01-01T00:00:00</ns3:Stop>" +
		"				<ns3:UserAccess>" +
		"					<ns3:CardExtension>32</ns3:CardExtension>" +
		"					<ns3:CardNo>" + t.AccessDevice.SmartcardCardNo + "</ns3:CardNo>" +
		"					<ns3:CardOrga>" + t.AccessDevice.SmartcardOrgaNo + "</ns3:CardOrga>" +
		"					<ns3:CocosSerialNo>0</ns3:CocosSerialNo>" +
		"					<ns3:PIN/>" +
		"					<ns3:PINs>0</ns3:PINs>" +
		"					<ns3:SerialNo>" + t.AccessDevice.SmartcardSerialNo + "</ns3:SerialNo>" +
		"					<ns3:TAN>0</ns3:TAN>" +
		"					<ns3:TempPIN/>" +
		"					<ns3:Type>" + t.AccessDevice.SmartcardType + "</ns3:Type>" +
		"				</ns3:UserAccess>" +
		"			</ns5:usage>" +
		"		</ns5:UsageEventReceived>" +
		"	</soap:Body>" +
		"</soap:Envelope>"
}

func generateProblemEvent(en EventName, t *Trip, ds *DriverSwipe) string {

	var (
		vehicleDevice     VehicleDevice
		smartcardType     string
		smartcardSerialNo string
		smartcardCardNo   string
		smartcardOrgaNo   string
		reservationId     string
		loc               *time.Location
	)

	if t != nil {
		reservationId = t.ReservationId
		vehicleDevice = t.VehicleDevice
		smartcardType = t.AccessDevice.SmartcardType
		smartcardSerialNo = t.AccessDevice.SmartcardSerialNo
		smartcardCardNo = t.AccessDevice.SmartcardCardNo
		smartcardOrgaNo = t.AccessDevice.SmartcardOrgaNo
		loc = t.Reservation.GetTimezone()

	} else if ds != nil {
		reservationId = "0"
		vehicleDevice = ds.VehicleDevice
		smartcardType = ds.AccessDevice.SmartcardType
		smartcardSerialNo = ds.AccessDevice.SmartcardSerialNo
		smartcardCardNo = ds.AccessDevice.SmartcardCardNo
		smartcardOrgaNo = ds.AccessDevice.SmartcardOrgaNo
		loc, _ = time.LoadLocation("UTC")
	}

	if smartcardType == "Hitag16" {
		smartcardType = "Hitag_16"
	} else if smartcardType == "Hitag32" {
		smartcardType = "Hitag_32"
	}

	return "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
		"	<soap:Body>" +
		"		<ns5:UsageProblemEventReceived xmlns=\"http://schemas.datacontract.org/2004/07/Invers.Ics.Interface\" xmlns:ns2=\"http://schemas.datacontract.org/2004/07/Invers.Ics.Interface.EvMo\" xmlns:ns3=\"http://invers.com\" xmlns:ns4=\"http://schemas.datacontract.org/2004/07/Invers.DataTypes\" xmlns:ns5=\"http://tempuri.org/\" xmlns:ns6=\"http://schemas.microsoft.com/2003/10/Serialization/\" xmlns:ns7=\"http://schemas.datacontract.org/2004/07/System.Net.Mail\">" +
		"			<ns5:usageProblem>" +
		"				<ns2:AdditionalParameters>" +
		"					<list>" +
		"						<AdditionalParameter>" +
		"							<Name>ReservationNo</Name>" +
		"							<ParamType>UInt32</ParamType>" +
		"							<Value>" + reservationId + "</Value>" +
		"						</AdditionalParameter>" +
		"					</list>" +
		"				</ns2:AdditionalParameters>" +
		"				<ns2:Description>" + fmt.Sprint(en) + "</ns2:Description>" +
		"				<ns2:Id>27642813</ns2:Id>" +
		"				<ns2:Position>" +
		"					<ns3:Altitude>0.0</ns3:Altitude>" +
		"					<ns3:Distance>0</ns3:Distance>" +
		"					<ns3:Format>ddd_dddddd</ns3:Format>" +
		"					<ns3:Latitude>51.49765166666667</ns3:Latitude>" +
		"					<ns3:LatitudeHemisphere>32</ns3:LatitudeHemisphere>" +
		"					<ns3:Longitude>-0.217075</ns3:Longitude>" +
		"					<ns3:LongitudeHemisphere>32</ns3:LongitudeHemisphere>" +
		"					<ns3:Quality>1</ns3:Quality>" +
		"					<ns3:SatInUse>8</ns3:SatInUse>" +
		"					<ns3:Timestamp>" + time.Now().In(loc).Format("2006-01-02T15:04:05") + "</ns3:Timestamp>" +
		"				</ns2:Position>" +
		"				<ns2:SentStatus>Sending</ns2:SentStatus>" +
		"				<ns2:Source>" +
		"					<ns3:DestinationAddress>" +
		"						<ns3:Fax/>" +
		"						<ns3:MailAddress>" +
		"							<ns3:BCC/>" +
		"							<ns3:CC/>" +
		"							<ns3:From>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns3:From>" +
		"							<ns3:Password/>" +
		"							<ns3:Priority>Normal</ns3:Priority>" +
		"							<ns3:ReplyTo>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns3:ReplyTo>" +
		"							<ns3:Sender>" +
		"								<Address/>" +
		"								<DisplayName/>" +
		"							</ns3:Sender>" +
		"							<ns3:Server/>" +
		"							<ns3:To/>" +
		"							<ns3:UserName/>" +
		"						</ns3:MailAddress>" +
		"						<ns3:PhoneNo>+" + vehicleDevice.VehiclePhoneNo + "</ns3:PhoneNo>" +
		"						<ns3:SIM>" +
		"							<GSMDataNo/>" +
		"							<GSMFaxNo/>" +
		"							<GSMProvider/>" +
		"							<GSMVoiceNo/>" +
		"							<ID/>" +
		"						</ns3:SIM>" +
		"						<ns3:TCPHost/>" +
		"						<ns3:TCPPort>0</ns3:TCPPort>" +
		"					</ns3:DestinationAddress>" +
		"					<ns3:DestinationType>BCSA</ns3:DestinationType>" +
		"					<ns3:Firmwareversion/>" +
		"					<ns3:OrgaNo>" + vehicleDevice.OrgaNo + "</ns3:OrgaNo>" +
		"					<ns3:SourceNo>132309508675338243</ns3:SourceNo>" +
		"				</ns2:Source>" +
		"				<ns2:SystemTimestamp xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:nil=\"true\"/>" +
		"				<ns2:Timestamp>" + time.Now().In(loc).Format("2006-01-02T15:04:05") + "</ns2:Timestamp>" +
		"				<ns2:Tlv/>" +
		"				<ns2:Type>12</ns2:Type>" +
		"				<ns3:AnswerList/>" +
		"				<ns3:BcStatus>WaitingForPIN</ns3:BcStatus>" +
		"				<ns3:CallReason>Unknown</ns3:CallReason>" +
		"				<ns3:CentralLockState>" +
		"					<NewOpen>false</NewOpen>" +
		"					<OpenCmd>false</OpenCmd>" +
		"					<Reason>Card</Reason>" +
		"				</ns3:CentralLockState>" +
		"				<ns3:DataFob>0</ns3:DataFob>" +
		"				<ns3:Driver>false</ns3:Driver>" +
		"				<ns3:DrivingDistance>0</ns3:DrivingDistance>" +
		"				<ns3:EnterPassengerCount>0</ns3:EnterPassengerCount>" +
		"				<ns3:Fuel>-1</ns3:Fuel>" +
		"				<ns3:FuelCard>0</ns3:FuelCard>" +
		"				<ns3:LedStatus>" +
		"					<Green>false</Green>" +
		"					<Red>false</Red>" +
		"					<Yellow>false</Yellow>" +
		"				</ns3:LedStatus>" +
		"				<ns3:Mileage>0</ns3:Mileage>" +
		"				<ns3:PassengerCount>0</ns3:PassengerCount>" +
		"				<ns3:Pause>false</ns3:Pause>" +
		"				<ns3:PinData>" +
		"					<PINs>1</PINs>" +
		"					<Result>OK</Result>" +
		"					<Tries>1</Tries>" +
		"				</ns3:PinData>" +
		"				<ns3:ReservationItem>" +
		"					<ID>0</ID>" +
		"					<Name/>" +
		"					<OrgaNo>" + vehicleDevice.OrgaNo + "</OrgaNo>" +
		"					<Type>BCSA</Type>" +
		"				</ns3:ReservationItem>" +
		"				<ns3:ReservationTypeId>-1</ns3:ReservationTypeId>" +
		"				<ns3:SpeedAlert>" +
		"					<Delay>PT0S</Delay>" +
		"					<Limit>0</Limit>" +
		"					<Speed>0</Speed>" +
		"				</ns3:SpeedAlert>" +
		"				<ns3:Start>1900-01-01T00:00:00</ns3:Start>" +
		"				<ns3:Stop>1900-01-01T00:00:00</ns3:Stop>" +
		"				<ns3:UserAccess>" +
		"					<ns3:CardExtension>32</ns3:CardExtension>" +
		"					<ns3:CardNo>" + smartcardCardNo + "</ns3:CardNo>" +
		"					<ns3:CardOrga>" + smartcardOrgaNo + "</ns3:CardOrga>" +
		"					<ns3:CocosSerialNo>0</ns3:CocosSerialNo>" +
		"					<ns3:PIN/>" +
		"					<ns3:PINs>0</ns3:PINs>" +
		"					<ns3:SerialNo>" + smartcardSerialNo + "</ns3:SerialNo>" +
		"					<ns3:TAN>0</ns3:TAN>" +
		"					<ns3:TempPIN/>" +
		"					<ns3:Type>" + smartcardType + "</ns3:Type>" +
		"				</ns3:UserAccess>" +
		"				<ns3:RejectedAccessReason>NoReservation</ns3:RejectedAccessReason>" +
		"			</ns5:usageProblem>" +
		"		</ns5:UsageProblemEventReceived>" +
		"	</soap:Body>" +
		"</soap:Envelope>"
}
