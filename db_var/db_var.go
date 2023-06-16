package db_var

import "time"

type MyQReservationStruct struct {
	Number      uint64 `json:"number"`
	OTAID       string `json:"ota_id"`
	BookingCode string `json:"booking_code"`
}

type DataCancelReservationStruct struct {
	BookingCode string `json:"booking_code"`
	OTAID       string `json:"ota_id"`
}

type DataInsertReservationStruct struct {
	BookingCode       string    `json:"booking_code"`
	OTAID             string    `json:"ota_id"`
	ArrivalTimeStr    string    `json:"arrival_time_str"`
	ArrivalDate       time.Time `json:"arrival_date"`
	DepartureDate     time.Time `json:"departure_date"`
	Adult             uint64    `json:"adult"`
	Child             uint64    `json:"child"`
	RoomTypeCode      string    `json:"room_type_code"`
	BedTypeCode       string    `json:"bed_type_code"`
	FullName          string    `json:"full_name"`
	Street            string    `json:"street"`
	City              string    `json:"city"`
	PostalCode        string    `json:"postal_code"`
	Phone1            string    `json:"phone1"`
	Email             string    `json:"email"`
	RoomRateAmountStr string    `json:"room_rate_amount_str"`
	RoomRateCode      string    `json:"room_rate_code"`
}

type DataUpdateReservationStruct struct {
	Vendor            string    `json:"vendor"`
	BookingCode       string    `json:"booking_code"`
	OTAID             string    `json:"ota_id"`
	ArrivalTimeStr    string    `json:"arrival_time_str"`
	ArrivalDate       time.Time `json:"arrival_date"`
	DepartureDate     time.Time `json:"departure_date"`
	Adult             uint64    `json:"adult"`
	Child             uint64    `json:"child"`
	RoomTypeCode      string    `json:"room_type_code"`
	BedTypeCode       string    `json:"bed_type_code"`
	FullName          string    `json:"full_name"`
	Street            string    `json:"street"`
	City              string    `json:"city"`
	PostalCode        string    `json:"postal_code"`
	Phone1            string    `json:"phone1"`
	Email             string    `json:"email"`
	RoomRateAmountStr string    `json:"room_rate_amount_str"`
	RoomRateCode      string    `json:"room_rate_code"`
}

