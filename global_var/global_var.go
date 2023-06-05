package global_var

type TReservationStatus struct {
	Canceled, Book, Modify string
}

var ReservationStatus = TReservationStatus{
	Canceled: "C",
	Book:     "B",
	Modify:   "M",
}
