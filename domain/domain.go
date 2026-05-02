package domain

import (
	"call_taxi_back_end/database"
	dataBaseType "call_taxi_back_end/database/types"
	"call_taxi_back_end/types"
	"call_taxi_back_end/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func BindingDriver(cus *types.Customer) (*dataBaseType.Order, error) {
	driver, err := database.Rdb.GetNearDriver(cus.Departure_Lat, cus.Departure_Lng)
	if err != nil {
		return nil, errors.New(database.NEAR_DRIVER_NOT_FOUND)
	}

	// only the one customer can change driver status and create order
	if err != nil {
		utils.PrintError(database.GET_LOCK_FAIL, err)
		return nil, err
	}

	driverData, err := database.PsqlDB.GetDriver(driver.Id)

	if err != nil {
		return nil, err
	}

	driverData.Status = database.BUSY_DRIVER
	err = database.PsqlDB.BindingDriver(driverData)

	if err != nil {
		return nil, err
	}

	newOrder := dataBaseType.Order{
		DriverId:         driver.Id,
		CusId:            cus.ID,
		Departure_Lat:    cus.Departure_Lat,
		Departure_Lng:    cus.Departure_Lng,
		Destination_Lat:  cus.Destination_Lat,
		Destination_Lng:  cus.Destination_Lng,
		Departure_Addr:   cus.Departure_Addr,
		Destination_Addr: cus.Destination_Addr,
	}

	order, err := database.PsqlDB.InsertOrder(&newOrder)
	if err != nil {
		return nil, err
	}

	err = database.Rdb.DriverBindOrder(driver)

	if err != nil {
		return nil, err
	}

	updateOrderCache(order)

	tcpDriver := &types.Driver{
		ID:                 driverData.ID,
		Name:               driverData.Name,
		Status:             database.BUSY_DRIVER,
		Route:              driverData.Route,
		Lat:                driver.Lat,
		Lng:                driver.Lng,
		CusDeparture_Lat:   cus.Departure_Lat,
		CusDeparture_Lng:   cus.Departure_Lng,
		CusDestination_Lat: cus.Destination_Lat,
		CusDestination_Lng: cus.Destination_Lng,
		OrderId:            order.ID,
		CusId:              order.CusId,
	}

	err = updateFakeDriverStatus(tcpDriver)

	if err != nil {
		utils.PrintError("update fakeDriver error", err)
		return nil, err
	}

	return order, nil
}

func DriverWorking(driverId int) error {

	//check driver data is in the redis or not
	//if exist delete old data and input new data
	//in the availibale driver geoset and driver detail data

	driverLocation, _ := database.Rdb.GetDriver(driverId, database.AVAILABLE_DRIVER)

	if driverLocation != nil {
		database.Rdb.DeleteDriver(driverId, database.AVAILABLE_DRIVER)
	}

	driverLocation, _ = database.Rdb.GetDriver(driverId, database.BUSY_DRIVER)

	if driverLocation != nil {
		database.Rdb.DeleteDriver(driverId, database.BUSY_DRIVER)
	}

	driverDetial, _ := database.Rdb.GetDriverDetail(driverId)

	if driverDetial != nil {
		database.Rdb.DeleteDriverDetail(driverId)
	}

	driver, err := database.PsqlDB.GetDriver(driverId)

	if err != nil {
		return err
	}

	driver.Status = database.AVAILABLE_DRIVER
	err = database.PsqlDB.UpdateDriver(driver)

	if err != nil {
		return err
	}

	tcpDriver := &types.Driver{
		ID:     driver.ID,
		Name:   driver.Name,
		Status: driver.Status,
		Route:  driver.Route,
	}

	err = updateFakeDriverStatus(tcpDriver)

	if err != nil {
		utils.PrintError("update fakeDriver error", err)
		return err
	}

	driverPosition := &dataBaseType.Redis_DriverLocation{
		Id:  driver.ID,
		Lat: 0,
		Lng: 0,
	}

	driverDetail := &dataBaseType.Redis_DriverDetail{
		Id:    driver.ID,
		Name:  driver.Name,
		Route: driver.Route,
	}

	database.Rdb.InsertDriver(driverPosition, database.AVAILABLE_DRIVER)
	database.Rdb.InsertDriverDetail(driverDetail)
	return nil
}

func DriverOff(driverId int) error {
	driver, err := database.PsqlDB.GetDriver(driverId)

	if err != nil {
		return err
	}

	if driver.Status != database.AVAILABLE_DRIVER {
		return errors.New(database.BUSY_DRIVER)
	}

	driver.Status = database.DRIVER_OFF

	err = database.PsqlDB.UpdateDriver(driver)

	if err != nil {
		return err
	}

	tcpDriver := &types.Driver{
		ID:     driver.ID,
		Name:   driver.Name,
		Status: driver.Status,
		Route:  driver.Route,
	}

	err = updateFakeDriverStatus(tcpDriver)

	if err != nil {
		utils.PrintError("update fakeDriver error", err)
		return err
	}

	driverLocation, _ := database.Rdb.GetDriver(driverId, database.BUSY_DRIVER)

	if driverLocation != nil {
		database.Rdb.DeleteDriver(driverId, database.BUSY_DRIVER)
	}

	driverLocation, _ = database.Rdb.GetDriver(driverId, database.AVAILABLE_DRIVER)
	if driverLocation != nil {
		database.Rdb.DeleteDriver(driverId, database.AVAILABLE_DRIVER)
	}

	driverDetial, _ := database.Rdb.GetDriverDetail(driverId)
	if driverDetial != nil {
		database.Rdb.DeleteDriverDetail(driverId)
	}

	return nil
}

func OrderFinish(orderId int) error {
	orderData, err := database.PsqlDB.GetOrder(orderId)

	if err != nil {
		return err
	}

	orderData.CompleteTime = time.Now()
	orderData.Status = database.ORDER_FINISH
	order, err := database.PsqlDB.UpdateOrder(orderData)

	if err != nil {
		return err
	}

	driver, err := database.PsqlDB.GetDriver(orderData.DriverId)
	driver.Status = database.AVAILABLE_DRIVER
	err = database.PsqlDB.UpdateDriver(driver)

	if err != nil {
		return err
	}

	updateOrderCache(order)

	tcpDriver := &types.Driver{
		ID:                 order.DriverId,
		Status:             database.AVAILABLE_DRIVER,
		CusDeparture_Lat:   0.00,
		CusDeparture_Lng:   0.00,
		CusDestination_Lat: 0.00,
		CusDestination_Lng: 0.00,
		Task:               "",
		CusId:              0,
	}

	err = updateFakeDriverStatus(tcpDriver)

	if err != nil {
		utils.PrintError("update fakeDriver error", err)
		return err
	}

	return nil
}

func VerifySession(signedID string) (string, bool) {
	sessionKey := []byte(os.Getenv("SESSIONKEY"))
	parts := strings.Split(signedID, ":")
	if len(parts) != 2 {
		return "", false
	}
	uuid, sig := parts[0], parts[1]

	res, err := database.Rdb.CheckSessionKey(uuid)

	if err != nil {
		return "", false
	}

	if res == false {
		return "", false
	}

	h := hmac.New(sha256.New, sessionKey)
	h.Write([]byte(uuid))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return uuid, true

	}
	return "", false
}

func CancelOrder(orderId int) (*dataBaseType.Order, error) {
	uuid := uuid.New()
	lockId := "order:" + uuid.String()

	lock := database.Rdb.CreateLock(lockId)
	err := lock.Lock()

	if err != nil {
		utils.PrintError(database.GET_LOCK_FAIL, err)
		return nil, err
	}

	defer lock.Unlock()

	order, err := database.PsqlDB.GetOrder(orderId)

	if err != nil {
		return nil, err
	}

	if order.Status == database.ORDER_FINISH {
		return nil, errors.New(database.ORDER_FINISH)
	}

	order, err = database.PsqlDB.CancelOrder(orderId)

	if err != nil {
		return nil, err
	}

	driver, err := database.PsqlDB.GetDriver(order.DriverId)
	driver.Status = database.AVAILABLE_DRIVER
	err = database.PsqlDB.UpdateDriver(driver)

	if err != nil {
		return nil, err
	}

	updateOrderCache(order)

	tcpDriver := &types.Driver{
		ID:                 order.DriverId,
		Status:             database.AVAILABLE_DRIVER,
		CusDeparture_Lat:   0.00,
		CusDeparture_Lng:   0.00,
		CusDestination_Lat: 0.00,
		CusDestination_Lng: 0.00,
		Task:               "",
		CusId:              0,
	}

	err = updateFakeDriverStatus(tcpDriver)

	if err != nil {
		utils.PrintError("update fakeDriver error", err)
		return nil, err
	}

	return order, nil
}

func GetCustomer(sessionId string) int {
	cusId, _ := database.Rdb.GetCustomerId(sessionId)

	if cusId != 0 {
		return cusId
	}

	customer, err := database.PsqlDB.GetCustomer(sessionId)

	if err != nil {
		utils.PrintError("customer is not exist", err)
		return 0
	}

	return customer.ID

}

func GetLatestOrder(cusId int) (*types.Order, error) {
	var resOrder *types.Order
	uuid := uuid.New()
	lock := database.Rdb.CreateLock(uuid.String())
	err := lock.Lock()

	if err != nil {
		utils.PrintError(database.GET_LOCK_FAIL, err)
		return nil, err
	}

	defer lock.Unlock()

	qOrder, err := database.PsqlDB.GetLatestOrder(cusId)

	if err != nil {
		return nil, errors.New("order is not exist")
	}

	resOrder = &types.Order{
		ID:               qOrder.ID,
		Price:            qOrder.Price,
		Status:           qOrder.Status,
		CreateTime:       qOrder.CreateTime,
		CompleteTime:     qOrder.CompleteTime,
		DriverId:         qOrder.DriverId,
		CusId:            qOrder.CusId,
		Departure_Lat:    0,
		Departure_Lng:    0,
		Destination_Lat:  0,
		Destination_Lng:  0,
		Departure_Addr:   qOrder.Departure_Addr,
		Destination_Addr: qOrder.Destination_Addr,
	}
	return resOrder, nil
}

func updateOrderCache(order *dataBaseType.Order) error {
	nOrder := dataBaseType.Redis_Order{
		ID:               order.ID,
		Price:            order.Price,
		Status:           order.Status,
		CreateTime:       order.CreateTime,
		CompleteTime:     order.CompleteTime,
		DriverId:         order.DriverId,
		CusId:            order.CusId,
		Departure_Addr:   order.Departure_Addr,
		Destination_Addr: order.Destination_Addr,
	}

	database.Rdb.InsertOrder(&nOrder)
	return nil
}

func updateFakeDriverStatus(driver *types.Driver) error {
	data, _ := json.Marshal(driver)
	conn, err := net.DialTimeout("tcp", os.Getenv("FAKEDRIVERSERVER"), 5*time.Second)

	if err != nil {
		conn.Close()
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		utils.PrintError("", err)
		return err
	}
	conn.Close()

	return nil
}
