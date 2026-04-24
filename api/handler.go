package api

import (
	"call_taxi_back_end/database"
	databasetype "call_taxi_back_end/database/types"
	"call_taxi_back_end/domain"
	"call_taxi_back_end/types"
	"call_taxi_back_end/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// only for psql data base
func InsertDriver(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	driver := &databasetype.Driver{}
	err = convertData(driver, body, w)
	if err != nil {
		return
	}

	err = database.PsqlDB.InsertDriver(driver)

	if err != nil {
		utils.PrintError(INSERT_DRIVER_FAIL, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(SERVER_ERROR))
		return
	}
	w.Write([]byte(INSERT_DRIVER_SUCCESS))
}

// only for psql data base
func UpdateDriver(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	driver := &databasetype.Driver{}
	err = convertData(driver, body, w)
	if err != nil {
		return
	}

	err = database.PsqlDB.UpdateDriver(driver)

	if err != nil {
		utils.PrintError(UPDATE_DRIVER_FAIL, err)
		w.Write([]byte(SERVER_ERROR))
		return
	}
	w.Write([]byte(UPDATE_DRIVER_SUCCESS))
}

// only for psql data base
func DeleteDriver(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		utils.PrintError(CONVERT_DATA_FAIL, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(INVALID_DATA))
		return
	}

	err = database.PsqlDB.DeleteDriver(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(SERVER_ERROR))
		return
	}
}

// only for psql data base
func GetDriver(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		utils.PrintError(CONVERT_DATA_FAIL, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(INVALID_DATA))
		return
	}

	driver, err := database.PsqlDB.GetDriver(id)

	if err != nil {
		utils.PrintError(GET_DRIVER_DATA_FAIL, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(SERVER_ERROR))
		return
	}
	json.NewEncoder(w).Encode(driver)
}

// only for psql data base
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}
	order := &databasetype.Order{}
	err = convertData(order, body, w)

	if err != nil {
		return
	}
	_, err = database.PsqlDB.UpdateOrder(order)
	if err != nil {
		utils.PrintError(UPDATE_ORDER_FAIL, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(SERVER_ERROR))
		return
	}
	w.Write([]byte(UPDATE_ORDER_SUCCESS))
}

// only for psql data base
func GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		utils.PrintError(CONVERT_DATA_FAIL, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(INVALID_DATA))
		return
	}

	order, err := database.PsqlDB.GetOrder(id)

	if err != nil {
		utils.PrintError(GET_ORDER_DATA_FAIL, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(SERVER_ERROR))
		return
	}
	json.NewEncoder(w).Encode(order)
}

func CallTaxi(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	cus := &types.Customer{}
	err = convertData(cus, body, w)
	if err != nil {
		return
	}

	sessionId := getSessionId(r)
	cus.ID, err = database.Rdb.GetCustomerId(sessionId)

	if err != nil {
		utils.PrintError(SESSION_ID_NOT_EXIST, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(CUSTOMER_DATA_NOT_EXIST))
		return
	}

	order, err := domain.BindingDriver(cus)

	if err != nil {
		if err.Error() == NEAR_DRIVER_NOT_FOUND {
			message := types.Message{
				Message: NEAR_DRIVER_NOT_FOUND,
			}
			json.NewEncoder(w).Encode(message)
		} else {
			utils.PrintError(CALLTAXI_FAIl, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(SERVER_ERROR))
			return
		}
		return
	}
	rOrder := resOrder(order)
	json.NewEncoder(w).Encode(rOrder)
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	order := &types.Order{}
	err = convertData(order, body, w)
	if err != nil {
		return
	}
	fmt.Println(order)

	sessionId := getSessionId(r)
	_, err = database.Rdb.GetCustomerId(sessionId)

	if err != nil {
		utils.PrintError(SESSION_ID_NOT_EXIST, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(CUSTOMER_DATA_NOT_EXIST))
		return
	}

	cOrder, err := domain.CancelOrder(order.ID)

	if err != nil {
		if err.Error() == database.ORDER_FINISH {
			message := types.Message{
				Message: ORDER_FINISH,
			}
			json.NewEncoder(w).Encode(message)
		} else {
			utils.PrintError(CALLTAXI_FAIl, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(SERVER_ERROR))
		}

		return
	}
	rOrder := resOrder(cOrder)
	json.NewEncoder(w).Encode(rOrder)
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	var sessionID string
	var isValid bool

	cookie, err := r.Cookie(SESSION_KEY_NAME)
	if err == nil {
		var cusId int
		sessionID, isValid = domain.VerifySession(cookie.Value)

		if isValid {
			updateCoockie(w, cookie.Value)
			cusId = domain.GetCustomer(sessionID)
		}

		if cusId == 0 {
			isValid = false
		} else {
			json.NewEncoder(w).Encode(types.CustomerId{ID: cusId})
			return
		}
	}

	if err != nil || !isValid {
		fullSignedID := utils.CreateSignedID()
		sessionID = strings.Split(fullSignedID, ":")[0]
		customer, err := database.PsqlDB.InsertCustomer(fullSignedID)

		if err != nil {
			return
		}

		database.Rdb.InsertSessionKey(sessionID, customer.ID)

		http.SetCookie(w, &http.Cookie{
			Name:     SESSION_KEY_NAME,
			Value:    fullSignedID,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,  // 防止 JS 讀取
			Secure:   false, // 正式環境設為 true
		})
		json.NewEncoder(w).Encode(types.CustomerId{ID: customer.ID})
	}

}

func UpdateDriverLocation(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	driver := &types.Driver{}
	err = convertData(driver, body, w)
	if err != nil {
		return
	}

	if driver.Task == DRIVER_TASK_FINISH {
		err := domain.OrderFinish(driver.OrderId)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	workRoute := []databasetype.Redis_WorkRoute{}

	for _, wroute := range driver.WorkRoute {
		r := databasetype.Redis_WorkRoute{
			RouteString: wroute.RouteString,
			NodeStep:    wroute.NodeStep,
		}
		workRoute = append(workRoute, r)
	}

	redisDriver := &databasetype.Redis_DriverLocation{
		Id:        driver.ID,
		Lat:       driver.Lat,
		Lng:       driver.Lng,
		WorkRoute: workRoute,
		Status:    driver.Status,
		Task:      driver.Task,
		CusId:     driver.CusId,
		OrderId:   driver.OrderId,
		Count:     driver.Count,
	}

	err = database.Rdb.UpdateDriverLocation(redisDriver)

	if err != nil {
		utils.PrintError(UPDATE_DRIVER_LOCATION_FAIL, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(SERVER_ERROR))
		return
	}
}

func GetLatestOrder(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}

	cus := &types.CustomerId{}
	err = convertData(cus, body, w)
	if err != nil {
		return
	}

	order, err := domain.GetLatestOrder(cus.ID)

	if err != nil {
		message := types.Message{
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(message)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func DriverWork(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}
	driver := &types.Driver{}
	err = convertData(driver, body, w)
	if err != nil {
		return
	}
	err = domain.DriverWorking(driver.ID)

	if err != nil {
		utils.PrintError(DRIVER_WORKING_FAIL, err)
		return
	}
	w.Write([]byte(DRIVER_WORKING_SUCCESS))
}

func DriverOff(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(w, r)
	if err != nil {
		return
	}
	driver := &types.Driver{}
	err = convertData(driver, body, w)
	if err != nil {
		return
	}
	err = domain.DriverOff(driver.ID)

	if err != nil {
		utils.PrintError(DRIVER_WORKING_FAIL, err)
		return
	}
	w.Write([]byte(DRIVER_WORKING_SUCCESS))
}

func updateCoockie(w http.ResponseWriter, session string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_KEY_NAME,
		Value:    session,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
	})
}

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.PrintError(READ_REQUEST_DATA_FAIL, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(INVALID_DATA))
		return nil, err
	}

	if len(body) == 0 {
		w.Write([]byte(DATA_IS_EMPTY))
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New(BODY_IS_EMPTY)
		utils.PrintError("err", err)
		return nil, err
	}

	return body, nil
}

func convertData(object any, data []byte, w http.ResponseWriter) error {
	err := json.Unmarshal(data, object)
	if err != nil {
		utils.PrintError(CONVERT_DATA_FAIL, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(INVALID_DATA))
		return err
	}
	return nil
}

func getSessionId(r *http.Request) string {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == SESSION_KEY_NAME {
			return strings.Split(cookie.Value, ":")[0]
		}
	}
	return ""
}

func resOrder(order *databasetype.Order) *types.Order {
	return &types.Order{
		ID:               order.ID,
		Price:            order.Price,
		Status:           order.Status,
		CreateTime:       order.CreateTime,
		CompleteTime:     order.CompleteTime,
		DriverId:         order.DriverId,
		CusId:            order.CusId,
		Departure_Lat:    order.Departure_Lat,
		Departure_Lng:    order.Departure_Lng,
		Destination_Lat:  order.Destination_Lat,
		Destination_Lng:  order.Departure_Lng,
		Departure_Addr:   order.Departure_Addr,
		Destination_Addr: order.Destination_Addr,
	}
}

func createSchema(w http.ResponseWriter, r *http.Request) {
	err := database.PsqlDB.CreateSchema()

	if err != nil {
		log.Printf("create schema fail: %s \n", err)
		w.Write([]byte("create schema fail"))
		return
	}
	w.Write([]byte("create schema success"))

}
