package types

import "time"

type Redis_DriverLocation struct {
	Id        int               `redis:"id"`
	Lat       float64           `redis:"lat"`
	Lng       float64           `redis:"lng"`
	WorkRoute []Redis_WorkRoute `redis:"workroute"`
	Status    string            `redis:"status"`
	Task      string            `redis:"task"`
	CusId     int               `redis:"cusId"`
	OrderId   int               `redis:"orderId"`
	Count     int               `redis:"count"`
}

type Redis_DriverDetail struct {
	Id    int    `redis:"id"`
	Name  string `redis:"name"`
	Route int    `redis:"route"`
}

type Redis_Order struct {
	ID               int       `redis:"id"`
	Price            int       `redis:"price"`
	Status           string    `redis:"status"`
	CreateTime       time.Time `redis:"createTime"`
	CompleteTime     time.Time `redis:"completeTime"`
	DriverId         int       `redis:"driverId"`
	CusId            int       `redis:"cusId"`
	Departure_Addr   string    `redis:"departure_addr"`
	Destination_Addr string    `redis:"destination_addr"`
}

type Redis_Session struct {
	CreateTime time.Time `redis:"createTime"`
	CustomerId int       `redis:"customerId"`
	UserType   string    `redis:"userType"`
}

type Redis_WorkRoute struct {
	RouteString string `json:"routeString"`
	NodeStep    int    `json:"nodeStep"`
}
