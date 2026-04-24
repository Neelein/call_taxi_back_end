package types

import "time"

type Driver struct {
	ID     int    `pg:",pk,type:bigserial"`
	Name   string `pg:",notnull"`
	Route  int    `pg:",notnull"`
	Status string `pg:",notnull"`
}

type Order struct {
	ID               int       `pg:",pk"`
	Price            int       `pg:",notnull"`
	Status           string    `pg:",notnull"`
	CreateTime       time.Time `pg:",notnull"`
	CompleteTime     time.Time `pg:",notnull"`
	DriverId         int       `pg:",notnull"`
	CusId            int       `pg:",notnull"`
	Departure_Lat    float64   `pg:",notnull"`
	Departure_Lng    float64   `pg:",notnull"`
	Destination_Lat  float64   `pg:",notnull"`
	Destination_Lng  float64   `pg:",notnull"`
	Departure_Addr   string    `pg:",notnull"`
	Destination_Addr string    `pg:",notnull"`
}

type Customer struct {
	ID         int       `pg:",pk,type:bigserial"`
	SessionId  string    `pg:",notnull"`
	CreateTime time.Time `pg:",notnull"`
}
