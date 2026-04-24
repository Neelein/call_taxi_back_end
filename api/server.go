package api

import (
	"call_taxi_back_end/database"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	Addr string
}

func CreateServer() *Server {
	return &Server{
		Addr: os.Getenv("ADDR"),
	}
}

func (s *Server) Start() error {
	handler := createRouter()
	database.CreateRedisServerClient()
	database.CreatePsqlServerClient()
	defer database.Rdb.ConnectClose()
	defer database.PsqlDB.PsqlServerClientConnectClose()
	println("server start ", s.Addr)
	err := http.ListenAndServe(s.Addr, handler)
	if err != nil {
		return err
	}
	return nil
}

func createRouter() http.Handler {
	r := mux.NewRouter()
	api := r.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/createschema", createSchema).Methods(http.MethodGet)
	api.HandleFunc("/updatedriver", UpdateDriver).Methods(http.MethodPost)
	api.HandleFunc("/insertdriver", InsertDriver).Methods(http.MethodPost)
	api.HandleFunc("/deletedriver", DeleteDriver).Methods(http.MethodGet)
	api.HandleFunc("/getdriver", GetDriver).Methods(http.MethodGet)
	api.HandleFunc("/updateorder", UpdateOrder).Methods(http.MethodPost)
	api.HandleFunc("/getorder", GetOrder).Methods(http.MethodGet)
	api.HandleFunc("/calltaxi", CallTaxi).Methods(http.MethodPost)
	api.HandleFunc("/getsession", SessionHandler).Methods(http.MethodGet)
	api.HandleFunc("/updatedriverlocation", UpdateDriverLocation).Methods(http.MethodPost)
	api.HandleFunc("/driverwork", DriverWork).Methods(http.MethodPost)
	api.HandleFunc("/driveroff", DriverOff).Methods(http.MethodPost)
	api.HandleFunc("/cancelorder", CancelOrder).Methods(http.MethodPost)
	api.HandleFunc("/getlatestorder", GetLatestOrder).Methods(http.MethodPost)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500", "http://127.0.0.1", "https://neeleindev.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:            true,
		AllowCredentials: true,
	})
	handler := c.Handler(api)
	return handler
}
