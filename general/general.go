package general

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Document struct {
	Title   string
	URL     string
	Content struct {
		Articles []struct {
			Title      string
			URL        string
			Categories []string
			Info       string
		}
	}
}

var IP = "http://192.168.1.58"
var Port = "9000"

func SetHeader(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	// TODO Token
	req.Header.Set("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc1MDM2MzAsInJlZnJlc2giOmZhbHNlLCJ1c2VyIjoiU1lTVEVNIn0.fqA0wHYfmtxhfD13M7zBXOxVan0OxeWY0elzAvGKTxk")
	return req
}

func GetBedTypeCode(InputString string) (string, string) {
	output := strings.Split(InputString, "#")
	return output[0], output[1]
}

func GenerateToken() string {
	uintValLeng1 := SetRandomInteger()
	uintValLeng2 := SetRandomInteger()
	uintValLeng3 := SetRandomInteger()
	uintValLeng4 := SetRandomInteger()
	uintValLeng5 := SetRandomInteger()
	uintValLeng6 := SetRandomInteger()
	uintValLeng7 := SetRandomInteger()

	hexValue1 := strconv.FormatUint(uint64(uintValLeng1), 16)
	hexValue2 := strconv.FormatUint(uint64(uintValLeng2), 16)
	hexValue3 := strconv.FormatUint(uint64(uintValLeng3), 16)
	hexValue4 := strconv.FormatUint(uint64(uintValLeng4), 16)
	hexValue5 := strconv.FormatUint(uint64(uintValLeng5), 16)
	hexValue6 := strconv.FormatUint(uint64(uintValLeng6), 16)
	hexValue7 := strconv.FormatUint(uint64(uintValLeng7), 16)

	hexValue1 = SetLength8(hexValue1)
	hexValue2 = SetLength4(hexValue2)
	hexValue3 = SetLength4(hexValue3)
	hexValue4 = SetLength4(hexValue4)
	hexValue5 = SetLength4(hexValue5)
	hexValue6 = SetLength8(hexValue6)
	hexValue7 = SetLength4(hexValue7)

	upperHexStr1 := strings.ToUpper(hexValue1)
	upperHexStr2 := strings.ToUpper(hexValue2)
	upperHexStr3 := strings.ToUpper(hexValue3)
	upperHexStr4 := strings.ToUpper(hexValue4)
	upperHexStr5 := strings.ToUpper(hexValue5)
	upperHexStr6 := strings.ToUpper(hexValue6)
	upperHexStr7 := strings.ToUpper(hexValue7)

	TokenToOTA := upperHexStr1 + "-" + upperHexStr2 + "-" + upperHexStr3 + "-" + upperHexStr4 + "-" + upperHexStr5 + "-" + upperHexStr6 + "-" + upperHexStr7
	return TokenToOTA
}

func SetRandomInteger() uint {
	maxUint := uint(^uint(0))
	str := strconv.FormatUint(uint64(maxUint), 10)
	letterBytes := str
	buffer := make([]byte, 16)
	_, _ = rand.Read(buffer)
	otpCharsLength := len(letterBytes)
	for i := 0; i < 16; i++ {
		buffer[i] = letterBytes[int(buffer[i])%otpCharsLength]
	}

	randomValue := string(buffer)
	num, err := strconv.ParseUint(randomValue, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to uint:", err)
	}
	return uint(num)
}

func SetLength4(data string) string {
	maxLength4 := 4
	dataStr := data
	if len(dataStr) > maxLength4 {
		dataStr = dataStr[:maxLength4]
	}
	return dataStr
}

func SetLength8(data string) string {
	maxLength8 := 8
	dataStr := data
	if len(dataStr) > maxLength8 {
		dataStr = dataStr[:maxLength8]
	}
	return dataStr
}

func GetTimeStampX(ADateTime time.Time) string {
	// format time yyyy-mm-ddT00:00:00+08:00
	Result := time.Date(ADateTime.Year(), ADateTime.Month(), ADateTime.Day(), 0, 0, 0, 0, time.Local)
	dateFormat := "2006-01-02T15:04:05+08:00"
	dateString := Result.Format(dateFormat)
	return dateString
}

// Header of XML
func GenerateXMLString(XMLSendBody string) string {
	CCMSSMUser := "CCMSSMUser"
	CCMSSMPassword := "CCMSSMPassword"

	Result := "<soap:Envelope xmlns:soap='http://schemas.xmlsoap.org/soap/envelope/'>" +
	"<soap:Header>" +
	"<wsse:Security soap:mustUnderstand='1' xmlns:wsse='http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd'>" +
	"<wsse:UsernameToken>" +
	"<wsse:Username>" + CCMSSMUser + "</wsse:Username>" +
	"<wsse:Password Type='http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText'>" + CCMSSMPassword + "</wsse:Password>" +
	"</wsse:UsernameToken>" +
	"</wsse:Security>" +
	"</soap:Header>" +
	"<soap:Body>" +
	XMLSendBody +
	"</soap:Body>" +
	"</soap:Envelope>"

	filePath := "output.xml"

	err := ioutil.WriteFile(filePath, []byte(Result), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}

	fmt.Println("File created successfully:", filePath)

	// data := Document{}
	// err := json.Unmarshal([]byte(jsonString), &data)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// doc := etree.NewDocument()
	// // doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	// err := doc.WriteToFile("output.xml")
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	return Result
}

func SendXMLOTA_NotifReportRQ(BookingCode string, OTAID string, RestStatus string) {
	// TODO
	TimeStampX := GetTimeStampX(time.Now())
	EchoToken := GenerateToken()

	BodyXML := "<OTA_NotifReportRQ xmlns=='http://www.opentravel.org/OTA/2003/05' Version=='1.0' TimeStamp=='" + TimeStampX + "' EchoToken=='" + EchoToken + "'>" +
		"<Success/>" +
		"<NotifDetails>" +
		"<HotelNotifReport>" +
		"<HotelReservations>" +
		"<HotelReservation CreateDateTime=='2010-01-01T12:00:00' ResStatus=='" + RestStatus + "'>" +
		"<UniqueID Type='16' ID='" + OTAID + "'/>" +
		"<ResGlobalInfo>" +
		"<HotelReservationIDs>" +
		"<HotelReservationID ResID_Type='14' ResID_Value='" + BookingCode + "'/>" +
		"</HotelReservationIDs>" +
		"</ResGlobalInfo>" +
		"</HotelReservation>" +
		"</HotelReservations>" +
		"</HotelNotifReport>" +
		"</NotifDetails>" +
		"</OTA_NotifReportRQ>"

	GenerateXMLString(BodyXML)
}

// Avail
// func SendXMLOTA_HotelAvailNotifRQ(TimeStampX string, EchoToken string, AvailabilityCount string, StartX string, EndX string, RoomTypeCode string) string {
// 	Result := "<OTA_HotelAvailNotifRQ xmlns=='http://www.opentravel.org/OTA/2003/05' Version=='1.0' TimeStamp=='" + TimeStampX + "' EchoToken=='" + EchoToken + "'>" +
// 		"<POS>" +
// 		"<Source>" +
// 		"<RequestorID Type='22' ID='" + ProgramConfiguration.CCMSSMRequestorID + "'/>" +
// 		"</Source>" +
// 		"</POS>" +
// 		"<AvailStatusMessages HotelCode='" + ProgramConfiguration.CCMSSMHotelCode + "'>" +
// 		"<AvailStatusMessage BookingLimit='" + AvailabilityCount + "'>" +
// 		"<StatusApplicationControl Start='" + StartX + "' End='" + EndX + "' InvTypeCode='" + RoomTypeCode + "'/>" +
// 		"</AvailStatusMessage>" +
// 		"</AvailStatusMessages>" +
// 		"</OTA_HotelAvailNotifRQ>"
// 	return Result
// }

// Rate
func SendXMLOTA_HotelRateAmountNotifRQsBNL2() {

}
