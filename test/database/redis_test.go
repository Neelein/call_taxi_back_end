package database_test

import (
	"call_taxi_back_end/database"
	"call_taxi_back_end/database/types"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestRedisGetDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	Id := 1

	driver, err := database.Rdb.GetDriver(Id, database.AVAILABLE_DRIVER)

	if err != nil {
		t.Fatalf("read driver position from redis fail: %s", err)
	}
	t.Log(driver)
}

func TestRedisInsertDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	driver := types.Redis_DriverLocation{
		Id:  3,
		Lat: 0.0001,
		Lng: 0.0001,
	}

	err := database.Rdb.InsertDriver(&driver, database.AVAILABLE_DRIVER)

	if err != nil {
		t.Fatalf("write driver postion into redis fail : %s", err)
	}
}

func TestRedisInsertDriverDetail(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	driver := types.Redis_DriverDetail{
		Id:    1,
		Name:  "Jack",
		Route: 1,
	}

	err := database.Rdb.InsertDriverDetail(&driver)

	if err != nil {
		t.Fatalf("write driver detail into redis fail : %s", err)
	}
}

func TestRedisDeleteDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	err := database.Rdb.DeleteDriver(1, database.AVAILABLE_DRIVER)

	if err != nil {
		t.Fatalf("delete driver postion from redis fail : %s", err)
	}
}

func TestRedisGetDriverDetail(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	driverDetail, err := database.Rdb.GetDriverDetail(1)

	if err != nil {
		t.Fatalf("get driver detail data from redis fail : %s", err)
	}
	t.Log(driverDetail)
}

func TestRedisDeleteDriverDetail(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	err := database.Rdb.DeleteDriverDetail(1)

	if err != nil {
		t.Fatalf("delete driver detail data from redis fail : %s", err)
	}
}

func TestRedisDriverBindOrder(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	driver := types.Redis_DriverLocation{
		Id:  1,
		Lat: 0.0001,
		Lng: 0.0001,
	}

	err := database.Rdb.DriverBindOrder(&driver)

	if err != nil {
		t.Fatalf("binding driver fail : %s", err)
	}
}

func TestRedisInsertOrder(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	order := types.Redis_Order{
		ID:           1,
		Price:        100,
		Status:       "create",
		CreateTime:   time.Now(),
		CompleteTime: time.Now(),
		DriverId:     1,
		CusId:        1,
	}

	err := database.Rdb.InsertOrder(&order)

	if err != nil {
		t.Fatalf("Insert Order to redis fail : %s", err)
	}
}

func TestRedisGetNearDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	driverLocation, err := database.Rdb.GetNearDriver(0.001, 0.001)

	if err != nil {
		t.Fatalf("get near driver fail : %s", err)
	}
	t.Log(driverLocation)
}

func TestRedisInsertSessionKey(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	err := database.Rdb.InsertSessionKey("testing", 1)

	if err != nil {
		t.Fatalf("insert session key error : %s", err)
	}
}

func TestGetCustomerId(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	id, err := database.Rdb.GetCustomerId("testing")

	if err != nil {
		t.Fatalf("get customer id error : %s", err)
	}

	t.Log(id)
}

func TestUpdateDriverLocation(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreateRedisServerClient()

	driver := types.Redis_DriverLocation{
		Id:  3,
		Lat: 0.0001,
		Lng: 0.0001,
	}

	err := database.Rdb.UpdateDriverLocation(&driver)

	if err != nil {
		t.Fatalf("update driver location error : %s", err)
	}
}
