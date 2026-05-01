package database

import (
	"call_taxi_back_end/database/types"
	"call_taxi_back_end/utils"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Psql struct {
	db *pg.DB
}

var PsqlDB *Psql

func CreatePsqlServerClient() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("PSQLUSER"),
		Password: os.Getenv("PSQLPASSWORD"),
		Database: os.Getenv("PSQLDATABASE"),
		Addr:     os.Getenv("PSQLADDR"),
	})
	PsqlDB = &Psql{
		db: db,
	}
}

func (p Psql) PsqlServerClientConnectClose() {
	PsqlDB.db.Close()
}

// create data table in psql
func (p Psql) CreateSchema() error {
	models := []interface{}{
		// (*types.Driver)(nil),
		(*types.Customer)(nil),
	}

	for _, model := range models {
		err := PsqlDB.db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			log.Println("err")
			return err
		}
	}

	return nil
}

func (p Psql) InsertDriver(driver *types.Driver) error {
	_, err := PsqlDB.db.Model(driver).Insert()
	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return err
	}
	return nil
}

func (p Psql) UpdateDriver(driver *types.Driver) error {
	driverData := &types.Driver{}
	err := PsqlDB.db.Model(driverData).Where("id = ?", driver.ID).Select()
	if err != pg.ErrNoRows && err != nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return err
	}

	//check data is exist or not
	if err == pg.ErrNoRows {
		utils.PrintError(DRIVER_NOT_EXIST, err)
		return err
	}

	_, err = PsqlDB.db.Model(driver).Where("id = ?", driver.ID).Update()

	if err != nil {
		utils.PrintError(UPDATE_DATA_FAIL, err)
		return err
	}

	return nil
}

func (p Psql) DeleteDriver(id int) error {
	driver := &types.Driver{}

	_, err := PsqlDB.db.Model(driver).Where("id = ?", id).Delete()

	if err != nil {
		utils.PrintError(DELETE_DATA_FAIL, err)
		return err
	}

	return nil
}

func (p Psql) GetDriver(id int) (*types.Driver, error) {
	driver := &types.Driver{}
	err := PsqlDB.db.Model(driver).Where("id = ?", id).Select()
	if err != nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return nil, err
	}

	return driver, nil
}

func (p Psql) InsertOrder(newOrder *types.Order) (*types.Order, error) {
	order := &types.Order{
		Price:            100,
		Status:           "process",
		CreateTime:       time.Now(),
		CompleteTime:     time.Now(),
		DriverId:         newOrder.DriverId,
		CusId:            newOrder.CusId,
		Departure_Lat:    newOrder.Departure_Lat,
		Departure_Lng:    newOrder.Departure_Lng,
		Destination_Lat:  newOrder.Destination_Lat,
		Destination_Lng:  newOrder.Destination_Lng,
		Departure_Addr:   newOrder.Departure_Addr,
		Destination_Addr: newOrder.Destination_Addr,
	}

	_, err := PsqlDB.db.Model(order).Returning("*").Insert()
	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return nil, err
	}

	return order, nil
}

func (p Psql) UpdateOrder(order *types.Order) (*types.Order, error) {
	_, err := PsqlDB.db.Model(order).Where("id = ?", order.ID).Returning("*").Update()

	if err != nil {
		utils.PrintError(UPDATE_DATA_FAIL, err)
		return order, err
	}

	return order, nil
}

func (p Psql) GetOrder(id int) (*types.Order, error) {
	order := &types.Order{}

	err := PsqlDB.db.Model(order).Where("id = ?", id).Select()

	if err != nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return nil, err
	}

	return order, nil
}

func (p Psql) InsertCustomer(sessionId string) (*types.Customer, error) {
	customer := &types.Customer{
		ID:         0,
		SessionId:  sessionId,
		CreateTime: time.Now(),
	}

	_, err := PsqlDB.db.Model(customer).Returning("*").Insert()

	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return nil, err
	}

	return customer, nil
}

func (p Psql) GetCustomer(sessionId string) (*types.Customer, error) {
	customer := &types.Customer{}

	err := PsqlDB.db.Model(customer).Where("session_id = ?", sessionId).Select()

	if err != nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return nil, err
	}

	return customer, nil
}

func (p Psql) CancelOrder(orderId int) (*types.Order, error) {
	order := types.Order{}

	_, err := PsqlDB.db.Model(&order).Where("id = ?", orderId).Set("status = ?", ORDER_CANCEL).Returning("*").Update()

	if err != nil {
		utils.PrintError(ORDER_CANCEL_FAIL, err)
		return nil, errors.New(ORDER_CANCEL_FAIL)
	}
	return &order, nil
}

func (p Psql) GetLatestOrder(cusId int) (*types.Order, error) {
	order := types.Order{}

	err := PsqlDB.db.Model(&order).Where("cus_id = ?", cusId).Order("create_time DESC").Limit(1).Select()

	if err != nil {
		utils.PrintError("test", err)
		return nil, err
	}
	return &order, nil
}

func (p Psql) BindingDriver(driver *types.Driver) error {
	driverData := &types.Driver{}
	err := PsqlDB.db.Model(driverData).Where("id = ? AND  status = ?", driver.ID, AVAILABLE_DRIVER).Select()

	if err == pg.ErrNoRows {
		return errors.New("no data exist")
	}

	_, err = PsqlDB.db.Model(driver).Where("id = ? AND  status = ?", driver.ID, AVAILABLE_DRIVER).Update()
	if err != nil {
		utils.PrintError(UPDATE_DATA_FAIL, err)
		return err
	}

	return nil
}
