package client_car_rental

type CarOrderCreateRequest struct {
	OrderName    string `json:"orderName"`
	Type         string `json:"type"`
	VehicleClass string `json:"vehicleClass"`
	HP           int    `json:"hp"`
	Seats        string `json:"seats"`
}

type CarOrderUpdateRequest struct {
	OrderID      string  `json:"orderId"`
	OrderName    string  `json:"orderName"`
	Type         string  `json:"type"`
	VehicleClass string  `json:"vehicleClass"`
	HP           int     `json:"hp"`
	Seats        string  `json:"seats"`
	Mileage      float64 `json:"mileage"`
}

type CarOrder struct {
	ID           string  `json:"id"`
	OrderName    string  `json:"orderName"`
	Type         string  `json:"type"`
	VehicleClass string  `json:"vehicleClass"`
	HP           int     `json:"hp"`
	Seats        string  `json:"seats"`
	Mileage      float64 `json:"mileage"`
}
