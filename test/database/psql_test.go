package database_test

import (
	"call_taxi_back_end/database"
	"call_taxi_back_end/database/types"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestPsqlInsertDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()

	driver := &types.Driver{
		ID:     0,
		Name:   "Ben",
		Route:  8,
		Status: "off",
	}

	err := database.PsqlDB.InsertDriver(driver)

	if err != nil {
		t.Fatalf("insert data fail: %s", err)
	}
}

func TestPsqlUpdateDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()

	driver := &types.Driver{
		ID:     8,
		Name:   "Binson",
		Route:  3,
		Status: database.DRIVER_OFF,
	}

	err := database.PsqlDB.UpdateDriver(driver)

	if err != nil {
		t.Fatalf("update data fail: %s", err)
	}
}

func TestPsqlDeleteDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()

	err := database.PsqlDB.DeleteDriver(12)

	if err != nil {
		t.Fatalf("delete data fail: %s", err)
	}
}

// func TestPsqlInsertOrder(t *testing.T) {
// 	godotenv.Load("../../dotenv.env")
// 	database.CreatePsqlServerClient()

// 	drivers, err := database.PsqlDB.InsertOrder(4, 1, 0.0001, 0.0001, 0.0005, 0.0005)

// 	if err != nil {
// 		t.Fatalf("insert order data fail: %s", err)
// 	}
// 	t.Log(drivers.ID)
// }

func TestPsqlUpdateOrder(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()

	order := &types.Order{
		ID:           1,
		Price:        100,
		Status:       "process",
		CreateTime:   time.Now(),
		CompleteTime: time.Now(),
		DriverId:     1,
		CusId:        4,
	}

	_, err := database.PsqlDB.UpdateOrder(order)

	if err != nil {
		t.Fatalf("update order data fail: %s", err)
	}
}

func TestPsqlGetOrder(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()

	order, err := database.PsqlDB.GetOrder(1)

	if err != nil {
		t.Fatalf("get order data fail: %s", err)
	}
	t.Log(order)
}

func TestPsqlInsertCustomer(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()

	customer, err := database.PsqlDB.InsertCustomer("testing2")

	if err != nil {
		t.Fatalf("insert customer data fail: %s", err)
	}
	t.Log(customer)
}

func TestPsqlGetDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	driver, err := database.PsqlDB.GetDriver(4)

	if err != nil {
		t.Fatalf("get driver fail: %s", err)
	}

	t.Log(driver)
}

func TestPsqlCancelOrder(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	order, err := database.PsqlDB.CancelOrder(100)

	if err != nil {
		t.Fatalf("get driver fail: %s", err)
	}

	t.Log(order)

}

func TestPsqlGetCustomer(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	cus, err := database.PsqlDB.GetCustomer("testing")
	if err != nil {
		t.Fatalf("get session fail: %s", err)
	}
	t.Log(cus)
}

func TestPsqlGetLatestOrder(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	cus, err := database.PsqlDB.GetLatestOrder(63)
	if err != nil {
		t.Fatalf("get session fail: %s", err)
	}
	t.Log(cus)
}

func TestPsqlBindingDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	driver := &types.Driver{
		ID:     7,
		Name:   "Binson",
		Route:  1,
		Status: database.BUSY_DRIVER,
	}

	err := database.PsqlDB.BindingDriver(driver)
	if err != nil {
		t.Fatalf("get session fail: %s", err)
	}
}
