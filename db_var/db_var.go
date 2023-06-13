package db_var

import ()

type MyQReservationStruct struct {
	Number      uint64 `json:"number"`
	OTAID       string `json:"ota_id"`
	BookingCode string `json:"booking_code"`
}

type DataReservationStatusStruct struct {
	ReservationNumber uint64 `json:"ReservationNumber"`
	StatusCode        string `json:"StatusCode"`
	CancelledBy       string `json:"CancelledBy"`
	CancelReason      string `json:"CancelReason"`
}

type DataReservationIsCMConfirmedStruct struct {
	BookingCode string `json:"BookingCode"`
	OTAID       string `json:"OTAID"`
	StatusCode  string `json:"StatusCode"`
	// TODO : value untuk siapa ini si "false"?
	False string `json:"False"`
}

