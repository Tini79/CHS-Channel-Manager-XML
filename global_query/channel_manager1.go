package global_query

import (
	"bytes"
	DBVar "channel-manager/db_var"
	General "channel-manager/general"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/beevik/etree"
)

func ChannelManager1ReadXML() {
	type QueryParamStruct struct {
		BookingCode string `json:"BookingCode"`
		OTAID       string `json:"OTAID"`
		HotelCode   string `json:"hotel_code"`
	}

	type DataXMLDetailStruct struct {
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

	type DataXMLStruct struct {
		ResStatus         string                `json:"ResStatus"`
		RoomRateCode      string                `json:"RoomRateCode"`
		AdultStr          string                `json:"AdultStr"`
		ChildStr          string                `json:"ChildStr"`
		InfantStr         string                `json:"InfantStr"`
		ArrivalDateStr    time.Time             `json:"ArrivalDateStr"`
		DepartureDateStr  time.Time             `json:"DepartureDateStr"`
		BookingCode       string                `json:"BookingCode"`
		OTAID             string                `json:"OTAID"`
		RoomRateAmountStr string                `json:"RoomRateAmountStr"`
		Details           []DataXMLDetailStruct `json:"details"`
	}

	var DataCancelReservation DBVar.DataCancelReservationStruct
	var DataInsertReservation DBVar.DataInsertReservationStruct
	var DataUpdateReservation DBVar.DataUpdateReservationStruct
	var DataXML DataXMLStruct
	var DataXMLDetail DataXMLDetailStruct
	var QueryParam QueryParamStruct

	doc := etree.NewDocument()
	// if err := doc.ReadFromFile("./file_xml/book.xml"); err != nil {
	if err := doc.ReadFromFile("./file_xml/reservation.xml"); err != nil {

		fmt.Println(err.Error())
	}

	root := doc.FindElement("//OTA_ResRetrieveRS")
	BedTypeCode := ""
	RoomTypeCode := ""
	DataXML.BookingCode = ""
	DataXML.OTAID = ""
	RPH := ""

	QueryParam.HotelCode = root.FindElement("//RoomStay/BasicPropertyInfo").SelectAttr("HotelCode").Value
	ReservationsList := root.FindElement("//ReservationsList")
	HotelReservation := ReservationsList.FindElements("HotelReservation")
	var ResStatus string

	for _, reservation := range HotelReservation {
		ResStatus = reservation.SelectAttr("ResStatus").Value
		CountHotelReservationField := reservation.FindElements("//UniqueID")
		// Unique ID
		for _, count := range CountHotelReservationField {
			if count.SelectAttr("Type").Value == "14" {
				DataXML.BookingCode = count.SelectAttr("ID").Value
				QueryParam.BookingCode = count.SelectAttr("ID").Value
			} else if count.SelectAttr("Type").Value == "16" {
				DataXML.OTAID = count.SelectAttr("ID").Value
				QueryParam.OTAID = count.SelectAttr("ID").Value
			}
		}

		RoomType := reservation.FindElement("//RoomType").SelectAttr("RoomTypeCode").Value
		RoomTypeCode, BedTypeCode = General.GetBedTypeCode(RoomType)
		RoomRate := reservation.FindElement("//RoomRate")
		DataXML.RoomRateCode = RoomRate.SelectAttr("RatePlanCode").Value
		DataXML.RoomRateAmountStr = reservation.FindElement("//Rate/Total").SelectAttr("AmountAfterTax").Value

		GuestCount := reservation.FindElements("//GuestCounts/GuestCount")
		for _, count := range GuestCount {
			if count.SelectAttr("AgeQualifyingCode").Value == "10" {
				DataXML.AdultStr = count.SelectAttr("Count").Value
			} else if count.SelectAttr("AgeQualifyingCode").Value == "8" {
				DataXML.ChildStr = count.SelectAttr("Count").Value
			} else if count.SelectAttr("AgeQualifyingCode").Value == "7" {
				DataXML.InfantStr = count.SelectAttr("Count").Value
			}
		}
		// Arrival Date & Departure Date
		layout := "2006-01-02"
		ArrivalDateStr, err := time.Parse(layout, reservation.FindElement("//TimeSpan").SelectAttr("Start").Value)
		if err != nil {
			fmt.Println(err)
		}
		DataXML.ArrivalDateStr = ArrivalDateStr

		DepartureDateStr, err := time.Parse(layout, reservation.FindElement("//TimeSpan").SelectAttr("End").Value)
		if err != nil {
			fmt.Println(err)
		}
		DataXML.DepartureDateStr = DepartureDateStr
		// tmbah dsini guest profile, tungguu file log dari pak krissss
		ResGuest := reservation.FindElements("//ResGuests/ResGuest")
		for _, count := range ResGuest {
			Profiles := count.SelectElement("Profiles")
			ProfileInfo := Profiles.SelectElement("ProfileInfo")
			Profile := ProfileInfo.SelectElement("Profile")
			Customer := Profile.SelectElement("Customer")
			PersonName := Customer.SelectElement(("PersonName"))
			Address := Customer.SelectElement(("Address"))
			DataXMLDetail.ResGuestRPH = count.SelectAttr("ResGuestRPH").Value
			DataXMLDetail.ArrivalTimeStr = count.SelectAttr("ArrivalTime").Value
			if PersonName.SelectElement("MiddleName") != nil {
				DataXMLDetail.MiddleName = PersonName.SelectElement("MiddleName").Text()
			}
			DataXMLDetail.GivenName = PersonName.SelectElement("GivenName").Text()
			DataXMLDetail.Surname = PersonName.SelectElement("Surname").Text()
			DataXMLDetail.Phone1 = count.FindElement("//Telephone").SelectAttr("PhoneNumber").Value
			DataXMLDetail.Email = Customer.SelectElement("Email").Text()
			DataXMLDetail.Street = Address.SelectElement("AddressLine").Text()
			DataXMLDetail.City = Address.SelectElement("CityName").Text()
			DataXMLDetail.PostalCode = Address.SelectElement("PostalCode").Text()
			DataXMLDetail.State = Address.SelectElement("StateProv").Text()
			DataXMLDetail.Country = Address.SelectElement("CountryName").Text()
			DataXMLDetail.Company = Address.SelectElement("CompanyName").Text()

			RPH = DataXMLDetail.ResGuestRPH

			// Insert data to Array of Details
			DataXML.Details = append(DataXML.Details, DataXMLDetail)
		}
	}
	log.Println(DataXML.Details, "bed")

	log.Println(RoomTypeCode, "typeeeee")

	MyQReservation := GetReservationByBookingCode(QueryParam.HotelCode, QueryParam.BookingCode, QueryParam.OTAID)

	if ResStatus == "Book" {
		if RPH != "" {
			if BedTypeCode == "" {
				// RoomList.Text := GetAvailableRoomByType(RoomTypeCode, '', FormatDateTimeX(ArrivalDate), FormatDateTimeX(DepartureDate), 0, 0, 0, 0, False, ProgramConfiguration.CCMSReservationAsAllotment);
				// if RoomList.Count > 0 {
				// BedTypeCode := GetBedTypeCode(RoomList.Strings[0]);
				// }
			}
			if RPH == DataXMLDetail.ResGuestRPH {
				Adult, err := strconv.ParseUint(DataXML.AdultStr, 10, 64)
				Child, err := strconv.ParseUint(DataXML.ChildStr, 10, 64)
				for _, detailsData := range DataXML.Details {
					DataInsertReservation.BookingCode = DataXML.BookingCode
					DataInsertReservation.OTAID = DataXML.OTAID
					DataInsertReservation.ArrivalDate = DataXML.ArrivalDateStr
					DataInsertReservation.DepartureDate = DataXML.DepartureDateStr
					DataInsertReservation.Adult = Adult
					DataInsertReservation.Child = Child
					DataInsertReservation.RoomTypeCode = RoomTypeCode
					// TODO BEDTYPE
					DataInsertReservation.BedTypeCode = BedTypeCode
					DataInsertReservation.ArrivalTimeStr = detailsData.ArrivalTimeStr
					DataInsertReservation.FullName = detailsData.GivenName + " " + detailsData.MiddleName + " " + detailsData.Surname
					log.Println(DataInsertReservation.FullName, "DataInsertReservation.FullName")
					DataInsertReservation.Street = detailsData.Street
					DataInsertReservation.City = detailsData.City
					DataInsertReservation.PostalCode = detailsData.PostalCode
					DataInsertReservation.Phone1 = detailsData.Phone1
					DataInsertReservation.Email = detailsData.Email
					DataInsertReservation.RoomRateAmountStr = DataXML.RoomRateAmountStr
					DataInsertReservation.RoomRateCode = DataXML.RoomRateCode
					InsertReservation(DataInsertReservation, QueryParam.HotelCode)
				}

				if err != nil {
					fmt.Println(err)
				}
			}
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
				if RPH == DataXMLDetail.ResGuestRPH {
					log.Println("uodate")
					Adult, err := strconv.ParseUint(DataXML.AdultStr, 10, 64)
					Child, err := strconv.ParseUint(DataXML.ChildStr, 10, 64)
					var DataUpdateReservationArr []DBVar.DataUpdateReservationStruct
					for _, detailsData := range DataXML.Details {
						DataUpdateReservation.BookingCode = DataXML.BookingCode
						DataUpdateReservation.OTAID = DataXML.OTAID
						DataUpdateReservation.ArrivalDate = DataXML.ArrivalDateStr
						DataUpdateReservation.DepartureDate = DataXML.DepartureDateStr
						DataUpdateReservation.Adult = Adult
						DataUpdateReservation.Child = Child
						DataUpdateReservation.RoomTypeCode = RoomTypeCode
						// TODO BEDTYPE
						DataUpdateReservation.BedTypeCode = BedTypeCode
						DataUpdateReservation.ArrivalTimeStr = detailsData.ArrivalTimeStr
						DataUpdateReservation.FullName = detailsData.GivenName + " " + detailsData.MiddleName + " " + detailsData.Surname
						log.Println(DataUpdateReservation.FullName, "fullll")
						DataUpdateReservation.Street = detailsData.Street
						DataUpdateReservation.City = detailsData.City
						DataUpdateReservation.PostalCode = detailsData.PostalCode
						DataUpdateReservation.Phone1 = detailsData.Phone1
						DataUpdateReservation.Email = detailsData.Email
						DataUpdateReservation.RoomRateAmountStr = DataXML.RoomRateAmountStr
						DataUpdateReservation.RoomRateCode = DataXML.RoomRateCode
						DataUpdateReservationArr = append(DataUpdateReservationArr, DataUpdateReservation)
					}
					UpdateReservation(DataUpdateReservationArr, QueryParam.HotelCode)

					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else if ResStatus == "Cancel" {
			if len(MyQReservation) != 0 {
				for _, detailsData := range MyQReservation {
					DataCancelReservation.BookingCode = detailsData.BookingCode
					DataCancelReservation.OTAID = detailsData.OTAID
					CancelReservation(DataCancelReservation, QueryParam.HotelCode)
				}
			}
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
	// TODO buat test, nanti hapus dan uncomment yg d bwh
	params.Set("OTAID", "")
	// params.Set("OTAID", OTAID)
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

func InsertReservation(DataXML DBVar.DataInsertReservationStruct, HotelCode string) {
	client := &http.Client{}
	payload, err := json.Marshal(DataXML)
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
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Insert request failed with status:", resp.StatusCode)
	}
}

func UpdateReservation(DataXML []DBVar.DataUpdateReservationStruct, HotelCode string) {
	client := &http.Client{}
	type DataToUpdateStruct struct {
		DataXML []DBVar.DataUpdateReservationStruct
	}

	var DataToUpdate DataToUpdateStruct
	DataToUpdate.DataXML = DataXML
	log.Println(DataToUpdate.DataXML, "DataToUpdate.DataXML")
	// DataToUpdate := {DataXML}
	payload, err := json.Marshal(DataXML)

	if err != nil {
		fmt.Println("Failed to create request:", err)
	}
	endPoint := fmt.Sprintf(General.IP + ":" + General.Port + "/UpdateReservation/" + HotelCode)
	req, err := http.NewRequest("PUT", endPoint, bytes.NewBuffer(payload))

	req = General.SetHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	log.Println(resp.Body, "payload")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Update request failed with status:", resp.StatusCode)
	}
}

func CancelReservation(DataXML DBVar.DataCancelReservationStruct, HotelCode string) {
	client := &http.Client{}
	payload, err := json.Marshal(DataXML)
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
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Cancel request failed with status:", resp.StatusCode)
	}
}
