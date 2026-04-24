package types

import (
	"time"

	"googlemaps.github.io/maps"
)

type Location struct {
	Departure_Lat   float64 `json:"departure_lat"`
	Departure_Lng   float64 `json:"departure_lng"`
	Destination_Lat float64 `json:"destination_lat"`
	Destination_Lng float64 `json:"destination_lng"`
}

type Route struct {
	Bounds   maps.LatLngBounds
	Polyline maps.Polyline
}

type DepatureLocation struct {
	Lat float64
	Lng float64
}

type WorkRoute struct {
	RouteString string `json:"routeString"`
	NodeStep    int    `json:"nodeStep"`
}

type Driver struct {
	ID                 int          `json:"id"`
	Name               string       `json:"name"`
	Route              int          `json:"route"`
	Status             string       `json:"status"`
	Lat                float64      `json:"lat"`
	Lng                float64      `json:"lng"`
	CusDeparture_Lat   float64      `json:"cusdeparture_lat"`
	CusDeparture_Lng   float64      `json:"cusdeparture_lng"`
	CusDestination_Lat float64      `json:"cusdestination_lat"`
	CusDestination_Lng float64      `json:"cusdestination_lng"`
	WorkRoute          []*WorkRoute `json:"workroute"`
	Count              int          `json:"count"`
	Task               string       `json:"task"`
	OrderId            int          `json:"orderId"`
	CusId              int          `json:"cusId"`
}

type RedisQuery struct {
	Key    string
	Member string
}

type Order struct {
	ID               int       `json:"id"`
	Price            int       `json:"price"`
	Status           string    `json:"status"`
	CreateTime       time.Time `json:"createTime"`
	CompleteTime     time.Time `json:"completeTime"`
	DriverId         int       `json:"driverId"`
	CusId            int       `json:"cusId"`
	Departure_Lat    float64   `json:"departure_lat"`
	Departure_Lng    float64   `json:"departure_lng"`
	Destination_Lat  float64   `json:"destination_lat"`
	Destination_Lng  float64   `json:"destination_lng"`
	Departure_Addr   string    `json:"departure_addr"`
	Destination_Addr string    `json:"destination_addr"`
}

type Customer struct {
	ID               int     `json:"id"`
	Departure_Lat    float64 `json:"departure_lat"`
	Departure_Lng    float64 `json:"departure_lng"`
	Destination_Lat  float64 `json:"destination_lat"`
	Destination_Lng  float64 `json:"destination_lng"`
	OrderId          int
	Departure_Addr   string `json:"departure_addr"`
	Destination_Addr string `json:"destination_addr"`
}

type CustomerId struct {
	ID int `json:"id"`
}

type Message struct {
	Message string `json:"message"`
}
