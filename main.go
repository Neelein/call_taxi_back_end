package main

import (
	"call_taxi_back_end/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("dotenv.env")
	server := api.CreateServer()
	log.Fatal(server.Start())
}
