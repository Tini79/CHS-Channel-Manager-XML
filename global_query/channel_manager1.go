package global_query

import (
	"bytes"
	GlobalVar "channel-manager/global_var"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/beevik/etree"
)

type Array map[string]interface{}
type MyQReservationStruct struct {
	Number      int    `json:"number"`
	OTAID       string `json:"ota_id"`
	BookingCode string `json:"booking_code"`
}
type ReservationStatusStruct struct {
	ReservationNumber int    `json:"ReservationNumber"`
	StatusCode        string `json:"StatusCode"`
	CancelledBy       string `json:"CancelledBy"`
	CancelReason      string `json:"CancelReason"`
}
type ReservationIsCMConfirmedStruct struct {
	BookingCode string `json:"BookingCode"`
	OTAID       string `json:"OTAID"`
	StatusCode  string `json:"StatusCode"`
	False       string `json:"False"`
}

var ReservationStatus ReservationStatusStruct
var ReservationIsCMConfirmed ReservationIsCMConfirmedStruct
var IP = "http://192.168.1.64:9000/"

func ChannelManager1ReadXML() {
	type QueryParamStruct struct {
		BookingCode string `json:"BookingCode"`
		OTAID       string `json:"OTAID"`
		HotelCode   string `json:"hotel_code"`
	}
	type DataInputDetailStruct struct {
		ResGuestRPH    string `json:"ResGuestRPH"`
		ArrivalTimeStr string `json:"ArrivalTimeStr"`
		GivenName      string `json:"GivenName"`
		MiddleName     string `json:"MiddleName"`
		Surname        string `json:"Surname"`
		Phone1         string `json:"Phone1"`
		Email          string `json:"Email"`
		Street         string `json:"Street"`
		City           string `json:"City"`
		PostalCode     string `json:"PostalCode"`
		State          string `json:"State"`
		Country        string `json:"Country"`
		Company        string `json:"Company"`
	}
	type DataInputStruct struct {
		ResStatus         string                  `json:"ResStatus"`
		RoomTypeCode      string                  `json:"RoomTypeCode"`
		RoomRateCode      string                  `json:"RoomRateCode"`
		AdultStr          string                  `json:"AdultStr"`
		ChildStr          string                  `json:"ChildStr"`
		InfantStr         string                  `json:"InfantStr"`
		ArrivalDateStr    time.Time               `json:"ArrivalDateStr"`
		DepartureDateStr  time.Time               `json:"DepartureDateStr"`
		BookingCode       string                  `json:"BookingCode"`
		OTAID             string                  `json:"OTAID"`
		RoomRateAmountStr string                  `json:"RoomRateAmountStr"`
		Details           []DataInputDetailStruct `json:"details"`
	}

	var DataInput DataInputStruct
	var DataInputDetail DataInputDetailStruct
	var QueryParam QueryParamStruct

	doc := etree.NewDocument()
	if err := doc.ReadFromFile("./book.xml"); err != nil {
		fmt.Println(err.Error())
	}

	root := doc.FindElement("//OTA_ResRetrieveRS")
	var BedTypeCode string
	var RPH string

	DataInput.BookingCode = ""
	DataInput.OTAID = ""

	// kirim hotel code juga
	QueryParam.HotelCode = root.FindElement("//RoomStay/BasicPropertyInfo").SelectAttr("HotelCode").Value
	ResStatus := root.FindElement("//HotelReservation").SelectAttr("ResStatus").Value

	// Unique ID
	CountHotelReservationField := root.FindElements("//UniqueID")
	for _, count := range CountHotelReservationField {
		if count.SelectAttr("Type").Value == "14" {
			DataInput.BookingCode = count.SelectAttr("ID").Value
			QueryParam.BookingCode = count.SelectAttr("ID").Value
		} else if count.SelectAttr("Type").Value == "16" {
			DataInput.OTAID = count.SelectAttr("ID").Value
			QueryParam.OTAID = count.SelectAttr("ID").Value
		}
	}

	// Room Stay
	// IsNoRPH := true

	// Room Rate
	RoomRate := root.FindElement("//RoomRate")
	// Room Rate Code
	DataInput.RoomRateCode = RoomRate.SelectAttr("RatePlanCode").Value
	// Room Type Code
	DataInput.RoomTypeCode = RoomRate.SelectAttr("RoomTypeCode").Value
	// Rate Amount After Tax
	DataInput.RoomRateAmountStr = root.FindElement("//Rate/Total").SelectAttr("AmountAfterTax").Value

	GuestCount := root.FindElements("//GuestCounts/GuestCount")
	for _, count := range GuestCount {
		if count.SelectAttr("AgeQualifyingCode").Value == "10" {
			// Adult
			DataInput.AdultStr = count.SelectAttr("Count").Value
		} else if count.SelectAttr("AgeQualifyingCode").Value == "8" {
			// Child
			DataInput.ChildStr = count.SelectAttr("Count").Value
		} else if count.SelectAttr("AgeQualifyingCode").Value == "7" {
			// Infant
			DataInput.InfantStr = count.SelectAttr("Count").Value
		}
	}

	// Arrival Date & Depature Date
	layout := "2006-01-02"
	ArrivalDateStr, err := time.Parse(layout, root.FindElement("//TimeSpan").SelectAttr("Start").Value)
	if err != nil {
		fmt.Println(err)
	}

	DataInput.ArrivalDateStr = ArrivalDateStr

	DepartureDateStr, err := time.Parse(layout, root.FindElement("//TimeSpan").SelectAttr("End").Value)
	if err != nil {
		fmt.Println(err)
	}
	DataInput.DepartureDateStr = DepartureDateStr

	// ResGuests
	ResGuests := root.FindElements("//ResGuests/ResGuest")
	for _, count := range ResGuests {
		Profiles := count.SelectElement("Profiles")
		ProfileInfo := Profiles.SelectElement("ProfileInfo")
		Profile := ProfileInfo.SelectElement("Profile")
		Customer := Profile.SelectElement("Customer")
		PersonName := Customer.SelectElement(("PersonName"))
		Address := Customer.SelectElement(("Address"))
		DataInputDetail.ResGuestRPH = count.SelectAttr("ResGuestRPH").Value
		DataInputDetail.ArrivalTimeStr = count.SelectAttr("ArrivalTime").Value
		if PersonName.SelectElement("MiddleName") != nil {
			DataInputDetail.MiddleName = PersonName.SelectElement("MiddleName").Text()
		}
		DataInputDetail.GivenName = PersonName.SelectElement("GivenName").Text()
		DataInputDetail.Surname = PersonName.SelectElement("Surname").Text()
		DataInputDetail.Phone1 = count.FindElement("//Telephone").SelectAttr("PhoneNumber").Value
		DataInputDetail.Email = Customer.SelectElement("Email").Text()
		DataInputDetail.Street = Address.SelectElement("AddressLine").Text()
		DataInputDetail.City = Address.SelectElement("CityName").Text()
		DataInputDetail.PostalCode = Address.SelectElement("PostalCode").Text()
		DataInputDetail.State = Address.SelectElement("StateProv").Text()
		DataInputDetail.Country = Address.SelectElement("CountryName").Text()
		DataInputDetail.Company = Address.SelectElement("CompanyName").Text()
		DataInput.Details = append(DataInput.Details, DataInputDetail)
	}

	MyQReservation := GetReservationByBookingCode(QueryParam.HotelCode, QueryParam.BookingCode, QueryParam.OTAID)

	if ResStatus == "Book" {
		if BedTypeCode == "" {
			// untuk mengisi bedtype get api dari pak khalil
			// RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
			// if RoomList.Count > 0 {
			// BedTypeCode := GetBedTypeCode(RoomList.Strings[0]);
			// }

			// for _, detailGuest := range DataInput.Details {
			if RPH == DataInputDetail.ResGuestRPH {
				// ServerDate := GetServerDate;
				// ProgramVariable.AuditDate := GetAuditDate;
				// ProgramConfiguration.CheckOutLimit := StrToTime(ReadConfigurationString(SystemCode.Hotel, ConfigurationCategory.Reservation, ConfigurationName.CheckOutLimit, False), ProgramVariable.FormatSettingX);
				// ReplaceTime(ArrivalDate, ArrivalTime);
				// ReplaceTime(DepartureDate, ProgramConfiguration.CheckOutLimit);

				// RoomRateAmount := 0.0
				// if (DataInput.RoomRateAmountStr != "") && (DataInput.RoomRateAmountStr != "0") {
				// 	RoomRateAmount, err = strconv.ParseFloat(DataInput.RoomRateAmountStr, 64)
				// 	if err != nil {
				// 		fmt.Println(err)
				// 	}
				// }
				// if (DateOf(ArrivalDate) >= DateOf(ProgramVariable.AuditDate)) and (DateOf(DepartureDate) > DateOf(ProgramVariable.AuditDate)) and (GetAvailableRoomCountByType(RoomTypeCode, BedTypeCode, ArrivalDate, DepartureDate, 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment) > 0) {
				// GuestProfileID := InsertGuestProfile('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
				//                                                    '', '', '', '', '', '', '', '', '', '', '', '',
				//                                                    '', '', '', '', '', '', '', '', '', '', '', '',
				//                                                    '', GuestProfileSource.Hotel, ServerDate);

				//               ContactPersonID := InsertContactPerson('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
				//                                                    '', '', '', '', '', '', '', '', '', '', '', '',
				//                                                    '', '', '', '', '', '', '', '', '', '', '', '');

				//               GuestDetailID := InsertGuestDetail(ArrivalDate, DepartureDate, Adult, Child, RoomTypeCode, BedTypeCode, '', RoomRateCode, '', '', '', '', '', '', True, RoomRateAmount, RoomRateAmount, 0, 0);
				//               ParameterCondition := InsertReservation(ContactPersonID, '', '', '', GuestDetailID, '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.New, '', '', '', '', '', '', BookingCode, OTAID, ResStatus, 1, NullDate, NullDate, True, False, ProgramConfiguration.CCMSReservationAsAllotment);
				//               AssignRoom(ParameterCondition, False);
				// } else {
				// 				GuestProfileID := InsertGuestProfile('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
				// 				'', '', '', '', '', '', '', '', '', '', '', '',
				// 				'', '', '', '', '', '', '', '', '', '', '', '',
				// 				'', GuestProfileSource.Hotel, ServerDate);

				// ContactPersonID := InsertContactPerson('', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00', CPType.Guest,
				// 				'', '', '', '', '', '', '', '', '', '', '', '',
				// 				'', '', '', '', '', '', '', '', '', '', '', '');

				// GuestDetailID := InsertGuestDetail(ArrivalDate, DepartureDate, Adult, Child, RoomTypeCode, BedTypeCode, '', RoomRateCode, '', '', '', '', '', '', True, RoomRateAmount, RoomRateAmount, 0, 0);
				// ParameterCondition := InsertReservation(ContactPersonID, '', '', '', GuestDetailID, '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.WaitList, '', '', '', '', '', '', BookingCode, OTAID, ResStatus, 1, NullDate, NullDate, True, False, ProgramConfiguration.CCMSReservationAsAllotment);
			}
			// }
		}

	} else {
		if ResStatus == "Modify" {
			if RPH != "" {
				if BedTypeCode == "" {
					// RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
					// if RoomList.Count > 0 {
					// BedTypeCode := GetBedTypeCode(RoomList.Strings[0]);
					// }
				}
				if RPH == DataInputDetail.ResGuestRPH {
					// 		if ProgramConfiguration.ChannelManagerVendor = ChannelManagerVendor.SiteMinder then
					// 		ChangeQueryString(MyQReservation,
					// 			'SELECT * FROM reservation' +
					// 			' WHERE booking_code="' +BookingCode+ '"' +
					// 			' AND booking_code<>"" ' +
					// 			'ORDER BY number;',
					// 			'', '', '', '', '', '', '', '', '', '')
					// 	else
					// 		ChangeQueryString(MyQReservation,
					// 			'SELECT * FROM reservation' +
					// 			' WHERE booking_code="' +BookingCode+ '"' +
					// 			' AND booking_code<>""' +
					// 			' AND ota_id="' +OTAID+ '"' +
					// 			' AND ota_id<>"" ' +
					// 			'ORDER BY number;',
					// 			'', '', '', '', '', '', '', '', '', '');
					// except

					// if len(MyQReservation) != 0 {
					// ProgramConfiguration.CheckOutLimit := StrToTime(ReadConfigurationString(SystemCode.Hotel, ConfigurationCategory.Reservation, ConfigurationName.CheckOutLimit, False), ProgramVariable.FormatSettingX);
					// ReplaceTime(ArrivalDate, ArrivalTime);
					// ReplaceTime(DepartureDate, ProgramConfiguration.CheckOutLimit);

					// RoomRateAmount := 0.0
					// if (DataInput.RoomRateAmountStr != "") && (DataInput.RoomRateAmountStr != "0") {
					// 	RoomRateAmount, err = strconv.ParseFloat(DataInput.RoomRateAmountStr, 64)
					// 	if err != nil {
					// 		fmt.Println(err)
					// 	}
					// }

					// for _, data := range MyQReservation {
					// fmt.Println(MyQReservation," data")
					// if data != nil {
					// ParameterCondition := MyQReservationnumber.AsLargeInt;
					// ContactPersonID := MyQReservationcontact_person_id.AsString;
					// GuestDetailID := MyQReservationguest_detail_id.AsString;
					// GuestProfileID := MyQReservationguest_profile_id.AsString;

					// UpdateContactPerson(ContactPersonID, '', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '', '0000-00-00',
					// 										'', '', '', '', '', '', '', '', '', '', '', '',
					// 										'', '', '', '', '', '', '', '', '', '', '', '');
					// UpdateGuestDetail(GuestDetailID, RoomTypeCode, BedTypeCode, '', RoomRateCode, '', '', '', '', '', '', ArrivalDate, DepartureDate, Adult, Child, True, RoomRateAmount, RoomRateAmount, 0, 0);
					// UpdateGuestProfile(GuestProfileID, '', FullName, Street, '', City, '', '', '', PostalCode, Phone1, '', '', Email, '', '', '', '', '', BoolToStr(True), '',  '0000-00-00',
					// 									 '', '', '', '', '', '', '', '', '', '', '', '',
					// 									 '', '', '', '', '', '', '', '', '', '', '', '',
					// 									 ServerDate);

					// 		if (DateOf(ArrivalDate) >= DateOf(ProgramVariable.AuditDate)) and (DateOf(DepartureDate) > DateOf(ProgramVariable.AuditDate)) and
					// 		(GetAvailableRoomCountByType(RoomTypeCode, BedTypeCode, ArrivalDate, DepartureDate, ParameterCondition, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment) > 0) then
					//  begin
					// 	 UpdateReservation(ParameterCondition, '', '', '', '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.New, '', '', '', '', '', '', OTAID, ResStatus, 1, NullDate, NullDate, True, False);

					// 	 AssignRoom(ParameterCondition, False)
					//  end
					//  else
					// 	 UpdateReservation(ParameterCondition, '', '', '', '', GuestProfileID, '', '', '', '', FullName, '', '', ReservationStatus.WaitList, '', '', '', '', '', '', OTAID, ResStatus, 1, NullDate, NullDate, True, False);
					// 	 //                            InsertLogUser(LogUserAction.InsertReservation, IntToStr(ParameterCondition), '', '', '', LogUserAction.InsertReservationX);
					// 	 //    //          ProcessSMSSchedule(SMSevent.OnInsertReservation, 'reservation.number = "' +IntToStr(ParameterCondition)+ '"', '', '', '', '', '', '', '');

					//  UpdateReservationIsCMConfirmed(MyQReservationbooking_code.AsString, MyQReservationota_id.AsString, OTAID, 'Modify', False);

					//  MyQReservation.Next;

					// }
					// }
					// }
				}
			}
		} else if ResStatus == "Cancel" {
			if len(MyQReservation) != 0 {
				for _, reservation := range MyQReservation {
					// UpdateReservationStatus
					ReservationStatus.ReservationNumber = reservation.Number
					ReservationStatus.StatusCode = GlobalVar.ReservationStatus.Canceled
					ReservationStatus.CancelledBy = ""
					ReservationStatus.CancelReason = "Cancel by Channel Manager"
					// UpdateReservationIsCMConfirmed
					ReservationIsCMConfirmed.BookingCode = reservation.BookingCode
					ReservationIsCMConfirmed.OTAID = reservation.OTAID
					ReservationIsCMConfirmed.StatusCode = GlobalVar.ReservationStatus.Canceled

					UpdateReservationStatusCancel(ReservationStatus, QueryParam.HotelCode)
					// // InsertLogUser(LogUserAction.CancelReservation, IntToStr(ReservationNumber), '', '', Reason, LogUserAction.CancelReservationX);

					// UpdateReservationIsCMConfirmed(reservation.BookingCode, reservation.OTAID, OTAID, 'Cancel', False);
				}
			}
		}
	}

}

func GetHeader(req *http.Request) *http.Request {
	// set Header
	req.Header.Set("Content-Type", "application/json")
	// TODO Token
	req.Header.Set("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc1MDM2MzAsInJlZnJlc2giOmZhbHNlLCJ1c2VyIjoiU1lTVEVNIn0.fqA0wHYfmtxhfD13M7zBXOxVan0OxeWY0elzAvGKTxk")

	return req
}

func GetReservationByBookingCode(HotelCode string, BookingCode string, OTAID string) []MyQReservationStruct {
	type RawDataStruct struct {
		Result json.RawMessage `json:"Result"`
	}

	var ReservationDataArray []MyQReservationStruct

	client := &http.Client{}
	endPoint := fmt.Sprintf(IP + "GetReservationByBookingCode/" + HotelCode)
	params := url.Values{}
	params.Set("BookingCode", BookingCode)
	// TODO buat test
	params.Set("OTAID", "")
	// params.Set("OTAID", OTAID)
	url := fmt.Sprintf("%s?%s", endPoint, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	req = GetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	var result RawDataStruct
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	err = json.Unmarshal(result.Result, &ReservationDataArray)
	if err != nil {
		fmt.Println(err)
	}

	return ReservationDataArray
}

func UpdateReservationStatusCancel(DataInput ReservationStatusStruct, HotelCode string) {
	client := &http.Client{}
	payload, err := json.Marshal(DataInput)
	if err != nil {
		fmt.Println("Failed to create request:", err)
	}
	endPoint := fmt.Sprintf(IP + "UpdateReservationStatus/" + HotelCode)
	req, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(payload))
	req = GetHeader(req)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Update request failed with status:", resp.StatusCode)
	}
}
