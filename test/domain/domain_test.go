package domain_test

import (
	"call_taxi_back_end/database"
	"call_taxi_back_end/domain"
	"call_taxi_back_end/types"
	"testing"

	"github.com/joho/godotenv"
)

func TestDriverWorking(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	database.CreateRedisServerClient()
	err := domain.DriverWorking(7)

	if err != nil {
		t.Fatalf("driver work fail: %s", err)
	}
}

func TestBindingDriver(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	database.CreateRedisServerClient()

	customer := types.Customer{
		ID:              1,
		Departure_Lat:   0.001,
		Departure_Lng:   0.001,
		Destination_Lat: 0.002,
		Destination_Lng: 0.002,
		OrderId:         0,
	}

	order, err := domain.BindingDriver(&customer)

	if err != nil {
		t.Fatalf("bindg driver fail: %s", err)
	}
	t.Log(order)
}

func TestDriverOff(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	database.CreateRedisServerClient()

	err := domain.DriverOff(4)

	if err != nil {
		t.Fatalf("driver off fail: %s", err)
	}
}

func TestOrderFinish(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	database.CreateRedisServerClient()

	err := domain.OrderFinish(6)

	if err != nil {
		t.Fatalf("complte order fail: %s", err)
	}
}

func TestVerifySession(t *testing.T) {
	godotenv.Load("../../dotenv.env")
	database.CreatePsqlServerClient()
	database.CreateRedisServerClient()

	testId := "94791d51-7988-42d9-9ed5-5b2a4c9ca2cb:3937c2e0988da61302a327a315eb2ae2894b05a0a2c4fbc6326c50ffb933fcd7"
	res, isvalid := domain.VerifySession(testId)

	t.Log(isvalid)
	t.Log(res)
}
