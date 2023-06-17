package global_query

import (
	"bytes"
	DBVar "channel-manager/db_var"
	General "channel-manager/general"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/beevik/etree"
)

var ParamMessageAPI string

func ChannelManager1ReadXML() {
	type QueryParamStruct struct {
		BookingCode string `json:"BookingCode"`
		OTAID       string `json:"OTAID"`
		HotelCode   string `json:"hotel_code"`
	}

	type DataResGuestXMLStruct struct {
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

	type DataHotelReservationXMLStruct struct {
		ResStatus         string
		BookingCode       string
		OTAID             string
		RoomType          string
		RoomTypeCode      string
		BedTypeCode       string
		RoomRateCode      string
		RoomRateAmountStr string
		AdultStr          string
		ChildStr          string
		InfantStr         string
		ArrivalDate       time.Time `json:"ArrivalDate"`
		DepartureDate     time.Time `json:"DepartureDate"`
	}

	type DataXMLStruct struct {
		// HotelReservationDetails []DataHotelReservationXMLStruct `json:""HotelReservationDetails`
		// Details                 []DataResGuestXMLStruct           `json:"details"`
	}

	var DataHotelReservationXML DataHotelReservationXMLStruct
	var DataResGuestXML DataResGuestXMLStruct
	var QueryParam QueryParamStruct
	var DataInsertReservation DBVar.DataInsertReservationStruct
	var DataInsertReservationArr []DBVar.DataInsertReservationStruct
	var DataUpdateReservation DBVar.DataUpdateReservationStruct
	var DataUpdateReservationArr []DBVar.DataUpdateReservationStruct
	var DataCancelReservation DBVar.DataCancelReservationStruct
	var DataCancelReservationArr []DBVar.DataCancelReservationStruct
	var DataInsertReservationIsCMConfirmed DBVar.DataReservationIsCMConfirmedStruct
	var DataInsertReservationIsCMConfirmedArr []DBVar.DataReservationIsCMConfirmedStruct
	var DataModifyReservationIsCMConfirmed DBVar.DataReservationIsCMConfirmedStruct
	var DataModifyReservationIsCMConfirmedArr []DBVar.DataReservationIsCMConfirmedStruct
	var DataCancelReservationIsCMConfirmed DBVar.DataReservationIsCMConfirmedStruct
	var DataCancelReservationIsCMConfirmedArr []DBVar.DataReservationIsCMConfirmedStruct

	doc := etree.NewDocument()
	// if err := doc.ReadFromFile("./file_xml/book.xml"); err != nil {
	if err := doc.ReadFromFile("./file_xml/reservation.xml"); err != nil {
		fmt.Println(err.Error())
	}
	root := doc.FindElement("//OTA_ResRetrieveRS")
	RPH := ""
	BookResStatus := ""
	ModifyResStatus := ""
	CancelResStatus := ""
	MyQReservation := GetReservationByBookingCode(QueryParam.HotelCode, QueryParam.BookingCode, QueryParam.OTAID)

	QueryParam.HotelCode = root.FindElement("//RoomStay/BasicPropertyInfo").SelectAttr("HotelCode").Value
	HotelReservations := root.FindElements("//ReservationsList/HotelReservation")
	// permasalahanya disinni loping 2x
	for _, reservation := range HotelReservations {
		DataHotelReservationXML.ResStatus = reservation.SelectAttr("ResStatus").Value
		CountHotelReservationField := reservation.SelectElements("UniqueID")
		// Unique ID
		for _, count := range CountHotelReservationField {
			if count.SelectAttr("Type").Value == "14" {
				DataHotelReservationXML.BookingCode = count.SelectAttr("ID").Value
				QueryParam.BookingCode = DataHotelReservationXML.BookingCode
			} else if count.SelectAttr("Type").Value == "16" {
				DataHotelReservationXML.OTAID = count.SelectAttr("ID").Value
				QueryParam.OTAID = DataHotelReservationXML.OTAID
			}
		}

		// Get RoomTypeCode and BedTypeCode
		RoomStays := reservation.SelectElement("RoomStays")
		RoomStay := RoomStays.SelectElement("RoomStay")
		RoomTypes := RoomStay.SelectElement("RoomTypes")
		DataHotelReservationXML.RoomType = RoomTypes.SelectElement("RoomType").SelectAttr("RoomTypeCode").Value
		DataHotelReservationXML.RoomTypeCode, DataHotelReservationXML.BedTypeCode = General.GetBedTypeCode(DataHotelReservationXML.RoomType)

		// Get RoomRate
		RoomRates := RoomStay.SelectElement("RoomRates")
		RoomRate := RoomRates.SelectElement("RoomRate")
		DataHotelReservationXML.RoomRateCode = RoomRate.SelectAttr("RatePlanCode").Value
		Rates := RoomRate.SelectElement("Rates")
		Rate := Rates.SelectElement("Rate")
		DataHotelReservationXML.RoomRateAmountStr = Rate.SelectElement("Total").SelectAttr("AmountAfterTax").Value

		// GuestCount
		GuestCounts := RoomStay.SelectElement("GuestCounts")
		GuestCount := GuestCounts.SelectElements("GuestCount")
		for _, count := range GuestCount {
			if count.SelectAttr("AgeQualifyingCode").Value == "10" {
				// Adult
				DataHotelReservationXML.AdultStr = count.SelectAttr("Count").Value
			} else if count.SelectAttr("AgeQualifyingCode").Value == "8" {
				// Child
				DataHotelReservationXML.ChildStr = count.SelectAttr("Count").Value
			} else if count.SelectAttr("AgeQualifyingCode").Value == "7" {
				// Infant
				DataHotelReservationXML.InfantStr = count.SelectAttr("Count").Value
			}
		}

		// Arrival Date
		layout := "2006-01-02"
		ArrivalDateStr, err := time.Parse(layout, RoomStay.SelectElement("TimeSpan").SelectAttr("Start").Value)
		if err != nil {
			fmt.Println(err)
		}
		DataHotelReservationXML.ArrivalDate = ArrivalDateStr

		// Departure Date
		DepartureDateStr, err := time.Parse(layout, RoomStay.SelectElement("TimeSpan").SelectAttr("End").Value)
		if err != nil {
			fmt.Println(err)
		}
		DataHotelReservationXML.DepartureDate = DepartureDateStr

		// Save Room Stay, #questionable
		if DataHotelReservationXML.RoomTypeCode != "" && DataHotelReservationXML.RoomRateCode != "" && !DataHotelReservationXML.ArrivalDate.IsZero() && !DataHotelReservationXML.DepartureDate.IsZero() {
			if DataHotelReservationXML.AdultStr == "" {
				DataHotelReservationXML.AdultStr = "1"
			}
			if DataHotelReservationXML.ChildStr == "" {
				DataHotelReservationXML.AdultStr = "0"
			}
			if DataHotelReservationXML.InfantStr == "" {
				DataHotelReservationXML.AdultStr = "0"
			}
		}

		// ResGuest
		ResGuests := reservation.SelectElement("ResGuests")
		ResGuest := ResGuests.SelectElement("ResGuest")
		DataResGuestXML.ResGuestRPH = ResGuest.SelectAttr("ResGuestRPH").Value
		RPH = DataResGuestXML.ResGuestRPH
		DataResGuestXML.ArrivalTimeStr = ResGuest.SelectAttr("ArrivalTime").Value
		Profiles := ResGuest.SelectElement("Profiles")
		ProfileInfo := Profiles.SelectElement("ProfileInfo")
		Profile := ProfileInfo.SelectElement("Profile")
		Customer := Profile.SelectElement("Customer")
		PersonName := Customer.SelectElement("PersonName")

		// Guest Personal Information
		DataResGuestXML.GivenName = PersonName.SelectElement("GivenName").Text()
		// DataResGuestXML.MiddleName = PersonName.SelectElement("MiddleName").Text()
		DataResGuestXML.Surname = PersonName.SelectElement("Surname").Text()
		DataResGuestXML.Phone1 = Customer.SelectElement("Telephone").SelectAttr("PhoneNumber").Value
		DataResGuestXML.Email = Customer.SelectElement("Email").Text()
		Address := Customer.SelectElement("Address")
		DataResGuestXML.Street = Address.SelectElement("AddressLine").Text()
		DataResGuestXML.City = Address.SelectElement("CityName").Text()
		DataResGuestXML.PostalCode = Address.SelectElement("PostalCode").Text()
		DataResGuestXML.State = Address.SelectElement("StateProv").Text()
		DataResGuestXML.Country = Address.SelectElement("CountryName").Text()
		DataResGuestXML.Company = Address.SelectElement("CompanyName").Text()

		if DataHotelReservationXML.ResStatus == "Book" {
			BookResStatus = DataHotelReservationXML.ResStatus
			if RPH != "" {
				if DataHotelReservationXML.BedTypeCode == "" {
					// RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
					// if RoomList.Count > 0 {
					// BedTypeCode := GetBedTypeCode(RoomList.Strings[0]);
					// }
				}
				if RPH == DataResGuestXML.ResGuestRPH {
					Adult, err := strconv.ParseUint(DataHotelReservationXML.AdultStr, 10, 64)
					if err != nil {
						fmt.Println(err)
					}
					Child, err := strconv.ParseUint(DataHotelReservationXML.ChildStr, 10, 64)
					if err != nil {
						fmt.Println(err)
					}

					// for _, detailsData := range DataXML.Details {
					// DataInsertReservation
					DataInsertReservation.BookingCode = DataHotelReservationXML.BookingCode
					DataInsertReservation.OTAID = DataHotelReservationXML.OTAID
					DataInsertReservation.ArrivalDate = DataHotelReservationXML.ArrivalDate
					DataInsertReservation.DepartureDate = DataHotelReservationXML.DepartureDate
					DataInsertReservation.Adult = Adult
					DataInsertReservation.Child = Child
					DataInsertReservation.RoomTypeCode = DataHotelReservationXML.RoomTypeCode
					DataInsertReservation.BedTypeCode = DataHotelReservationXML.BedTypeCode
					DataInsertReservation.ArrivalTimeStr = DataResGuestXML.ArrivalTimeStr
					DataInsertReservation.FullName = DataResGuestXML.GivenName + " " + DataResGuestXML.MiddleName + " " + DataResGuestXML.Surname
					DataInsertReservation.Street = DataResGuestXML.Street
					DataInsertReservation.City = DataResGuestXML.City
					DataInsertReservation.PostalCode = DataResGuestXML.PostalCode
					DataInsertReservation.Phone1 = DataResGuestXML.Phone1
					DataInsertReservation.Email = DataResGuestXML.Email
					DataInsertReservation.RoomRateAmountStr = DataHotelReservationXML.RoomRateAmountStr
					DataInsertReservation.RoomRateCode = DataHotelReservationXML.RoomRateCode
					DataInsertReservationArr = append(DataInsertReservationArr, DataInsertReservation)

					// UpdateReservationIsCMConfirmed
					DataInsertReservationIsCMConfirmed.BookingCode = DataHotelReservationXML.BookingCode
					DataInsertReservationIsCMConfirmed.OTAID = DataHotelReservationXML.OTAID
					DataInsertReservationIsCMConfirmed.IsCmConfirmed = true
					DataCancelReservationIsCMConfirmed.ResStatus = DataHotelReservationXML.ResStatus // Untuk dikirim ke NotifXML sebagai parameter
					DataInsertReservationIsCMConfirmedArr = append(DataInsertReservationIsCMConfirmedArr, DataInsertReservationIsCMConfirmed)
				}
				// InsertReservation(DataInsertReservationArr, QueryParam.HotelCode, DataReservationIsCMConfirmedArr)
			}
		} else {
			if DataHotelReservationXML.ResStatus == "Modify" {
				ModifyResStatus = DataHotelReservationXML.ResStatus
				if RPH != "" {
					if DataHotelReservationXML.BedTypeCode == "" {
						// RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
						// if RoomList.Count > 0 {
						// DataHotelReservationXML := GetDataHotelReservationXML(RoomList.Strings[0]);
						// }
					}
					if RPH == DataResGuestXML.ResGuestRPH {
						Adult, err := strconv.ParseUint(DataHotelReservationXML.AdultStr, 10, 64)
						if err != nil {
							fmt.Println(err)
						}
						Child, err := strconv.ParseUint(DataHotelReservationXML.ChildStr, 10, 64)
						if err != nil {
							fmt.Println(err)
						}

						// for _, detailsData := range DataXML.Details {
						DataUpdateReservation.BookingCode = DataHotelReservationXML.BookingCode
						DataUpdateReservation.OTAID = DataHotelReservationXML.OTAID
						DataUpdateReservation.ArrivalDate = DataHotelReservationXML.ArrivalDate
						DataUpdateReservation.DepartureDate = DataHotelReservationXML.DepartureDate
						DataUpdateReservation.Adult = Adult
						DataUpdateReservation.Child = Child
						DataUpdateReservation.RoomTypeCode = DataHotelReservationXML.RoomTypeCode
						DataUpdateReservation.BedTypeCode = DataHotelReservationXML.BedTypeCode
						DataUpdateReservation.ArrivalTimeStr = DataResGuestXML.ArrivalTimeStr
						DataUpdateReservation.FullName = DataResGuestXML.GivenName + " " + DataResGuestXML.MiddleName + " " + DataResGuestXML.Surname
						DataUpdateReservation.Street = DataResGuestXML.Street
						DataUpdateReservation.City = DataResGuestXML.City
						DataUpdateReservation.PostalCode = DataResGuestXML.PostalCode
						DataUpdateReservation.Phone1 = DataResGuestXML.Phone1
						DataUpdateReservation.Email = DataResGuestXML.Email
						DataUpdateReservation.RoomRateAmountStr = DataHotelReservationXML.RoomRateAmountStr
						DataUpdateReservation.RoomRateCode = DataHotelReservationXML.RoomRateCode
						DataUpdateReservationArr = append(DataUpdateReservationArr, DataUpdateReservation)

						// UpdateReservationIsCMConfirmed
						DataModifyReservationIsCMConfirmed.BookingCode = DataHotelReservationXML.BookingCode
						DataModifyReservationIsCMConfirmed.OTAID = DataHotelReservationXML.OTAID
						DataModifyReservationIsCMConfirmed.IsCmConfirmed = false
						DataCancelReservationIsCMConfirmed.ResStatus = DataHotelReservationXML.ResStatus // Untuk dikirim ke NotifXML sebagai parameter
						DataModifyReservationIsCMConfirmedArr = append(DataModifyReservationIsCMConfirmedArr, DataModifyReservationIsCMConfirmed)
						// }
						// UpdateReservation(DataUpdateReservationArr, QueryParam.HotelCode, DataReservationIsCMConfirmedArr)
					}
				}
			} else if DataHotelReservationXML.ResStatus == "Cancel" {
				CancelResStatus = DataHotelReservationXML.ResStatus
				if len(MyQReservation) != 0 {
					// for _, detailsData := range MyQReservation {
					DataCancelReservation.BookingCode = DataHotelReservationXML.BookingCode
					DataCancelReservation.OTAID = DataHotelReservationXML.OTAID

					// UpdateReservationIsCMConfirmed
					DataCancelReservationIsCMConfirmed.BookingCode = DataHotelReservationXML.BookingCode
					DataCancelReservationIsCMConfirmed.OTAID = DataHotelReservationXML.OTAID
					DataCancelReservationIsCMConfirmed.IsCmConfirmed = false
					DataCancelReservationIsCMConfirmed.ResStatus = DataHotelReservationXML.ResStatus // Untuk dikirim ke NotifXML sebagai parameter
					DataCancelReservationArr = append(DataCancelReservationArr, DataCancelReservation)
					DataCancelReservationIsCMConfirmedArr = append(DataCancelReservationIsCMConfirmedArr, DataCancelReservationIsCMConfirmed)
					// }
					// CancelReservation(DataCancelReservationArr, QueryParam.HotelCode, DataReservationIsCMConfirmedArr)
				}
			}
		}
	}

	if BookResStatus != "" {
		InsertReservation(DataInsertReservationArr, QueryParam.HotelCode, DataInsertReservationIsCMConfirmedArr)
	} else {
		if ModifyResStatus != "" {
			UpdateReservation(DataUpdateReservationArr, QueryParam.HotelCode, DataModifyReservationIsCMConfirmedArr)
		} else if CancelResStatus != "" {
			CancelReservation(DataCancelReservationArr, QueryParam.HotelCode, DataCancelReservationIsCMConfirmedArr)
		}
	}
}

func GetReservationByBookingCode(HotelCode string, BookingCode string, OTAID string) []DBVar.MyQReservationStruct {
	type RawDataStruct struct {
		Result json.RawMessage `json:"Result"`
	}

	var ReservationDataArray []DBVar.MyQReservationStruct

	client := &http.Client{}
	endPoint := fmt.Sprintf(General.IP + ":" + General.Port + "/GetReservationByBookingCode/" + HotelCode)
	params := url.Values{}
	params.Set("BookingCode", BookingCode)
	params.Set("OTAID", OTAID)
	url := fmt.Sprintf("%s?%s", endPoint, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	req = General.SetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

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

func InsertReservation(DataXML []DBVar.DataInsertReservationStruct, HotelCode string, DataReservationIsCMConfirmedArr []DBVar.DataReservationIsCMConfirmedStruct) {
	client := &http.Client{}
	type DataToInsertStruct struct {
		Data []DBVar.DataInsertReservationStruct `json:"data"`
	}

	var DataToInsert DataToInsertStruct
	DataToInsert.Data = DataXML
	log.Println(DataToInsert, "DataToInsert")
	payload, err := json.Marshal(DataToInsert)

	if err != nil {
		fmt.Println("Failed to create request:", err)
	}

	endPoint := fmt.Sprintf(General.IP + ":" + General.Port + "/InsertReservation/" + HotelCode)
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(payload))
	req = General.SetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
	if err != nil {
		return
	}

	UpdateReservationIsCMConfirmed(DataReservationIsCMConfirmedArr, HotelCode)

	// var DataReservationIsCMConfirmedArr2 []DBVar.DataReservationIsCMConfirmedStruct
	for _, DataReservation := range DataReservationIsCMConfirmedArr {
		General.SendXMLOTA_NotifReportRQ(DataReservation.BookingCode, DataReservation.OTAID, DataReservation.ResStatus)
		// DataReservation.IsCmConfirmed = true
		// DataReservationIsCMConfirmedArr2 = append(DataReservationIsCMConfirmedArr2, DataReservation)
	}
}

func UpdateReservation(DataXML []DBVar.DataUpdateReservationStruct, HotelCode string, DataReservationIsCMConfirmedArr []DBVar.DataReservationIsCMConfirmedStruct) {
	UpdateReservationIsCMConfirmed(DataReservationIsCMConfirmedArr, HotelCode)

	client := &http.Client{}
	type DataToUpdateStruct struct {
		Data []DBVar.DataUpdateReservationStruct `json:"data"`
	}

	var DataToUpdate DataToUpdateStruct
	DataToUpdate.Data = DataXML
	payload, err := json.Marshal(DataToUpdate)

	if err != nil {
		fmt.Println("Failed to create request:", err)
	}
	endPoint := fmt.Sprintf(General.IP + ":" + General.Port + "/UpdateReservation/" + HotelCode)
	req, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(payload))
	// TODO : TAMPILKAN RESPONS
	req = General.SetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
	if err != nil {
		fmt.Println(err)
		return
	}

	var DataReservationIsCMConfirmedArr2 []DBVar.DataReservationIsCMConfirmedStruct
	for _, DataReservation := range DataReservationIsCMConfirmedArr {
		General.SendXMLOTA_NotifReportRQ(DataReservation.BookingCode, DataReservation.OTAID, DataReservation.ResStatus)
		DataReservation.IsCmConfirmed = true
		DataReservationIsCMConfirmedArr2 = append(DataReservationIsCMConfirmedArr2, DataReservation)
	}
	// TODO jika error jangan jalankan func di bawah
	UpdateReservationIsCMConfirmed(DataReservationIsCMConfirmedArr2, HotelCode)
}

func CancelReservation(DataXML []DBVar.DataCancelReservationStruct, HotelCode string, DataReservationIsCMConfirmedArr []DBVar.DataReservationIsCMConfirmedStruct) {
	client := &http.Client{}
	type DataToCancelStruct struct {
		Data []DBVar.DataCancelReservationStruct `json:"data"`
	}

	var DataToCancel DataToCancelStruct
	DataToCancel.Data = DataXML
	payload, err := json.Marshal(DataToCancel)
	if err != nil {
		fmt.Println("Failed to create request:", err)
	}
	endPoint := fmt.Sprintf(General.IP + ":" + General.Port + "/CancelReservation/" + HotelCode)
	req, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(payload))
	req = General.SetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
	if err != nil {
		return
	}

	UpdateReservationIsCMConfirmed(DataReservationIsCMConfirmedArr, HotelCode)

	var DataReservationIsCMConfirmedArr2 []DBVar.DataReservationIsCMConfirmedStruct
	for _, DataReservation := range DataReservationIsCMConfirmedArr {
		General.SendXMLOTA_NotifReportRQ(DataReservation.BookingCode, DataReservation.OTAID, DataReservation.ResStatus)
		DataReservation.IsCmConfirmed = true
		DataReservationIsCMConfirmedArr2 = append(DataReservationIsCMConfirmedArr2, DataReservation)
	}

	UpdateReservationIsCMConfirmed(DataReservationIsCMConfirmedArr2, HotelCode)
}

func UpdateReservationIsCMConfirmed(DataXML []DBVar.DataReservationIsCMConfirmedStruct, HotelCode string) {
	client := &http.Client{}
	type DataToUpdateStruct struct {
		Data []DBVar.DataReservationIsCMConfirmedStruct `json:"data"`
	}

	var DataToUpdate DataToUpdateStruct
	DataToUpdate.Data = DataXML
	payload, err := json.Marshal(DataToUpdate)

	if err != nil {
		fmt.Println("Failed to create request:", err)
	}
	endPoint := fmt.Sprintf(General.IP + ":" + General.Port + "/UpdateReservationIsCMConfirmed/" + HotelCode)
	req, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(payload))
	// TODO : TAMPILKAN RESPONS
	req = General.SetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// deful 1

// update is confirm 0

// update reservation

// update is confirm 1
