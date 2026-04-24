package database

import (
	"call_taxi_back_end/database/types"
	"call_taxi_back_end/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	rs     *redsync.Redsync
}

var Rdb *Redis

var ctx = context.Background()

func CreateRedisServerClient() {
	DB, _ := strconv.Atoi(os.Getenv("REDISDB"))

	rdbClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISADDR"),
		Password: os.Getenv("REDISPASSWORD"),
		DB:       DB,
	})

	pool := goredis.NewPool(rdbClient)

	rs := redsync.New(pool)

	Rdb = &Redis{
		client: rdbClient,
		rs:     rs,
	}
}

func (r Redis) ConnectClose() {
	Rdb.client.Close()
}

func (r Redis) CreateLock(lockId string) *redsync.Mutex {
	mutex := Rdb.rs.NewMutex(lockId)
	return mutex
}

func (r Redis) InsertDriver(driver *types.Redis_DriverLocation, zsetName string) error {
	_, err := Rdb.client.GeoAdd(ctx, zsetName,
		&redis.GeoLocation{
			Name:      strconv.Itoa(driver.Id),
			Longitude: driver.Lng,
			Latitude:  driver.Lat,
		}).Result()

	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return err
	}

	return nil
}

func (r Redis) InsertDriverDetail(driver *types.Redis_DriverDetail) error {
	_, err := Rdb.client.HSet(ctx, DRIVER_DETAIL+":"+strconv.Itoa(driver.Id), driver).Result()

	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return err
	}

	return nil
}

func (r Redis) GetDriver(id int, zsetName string) (*types.Redis_DriverLocation, error) {
	dirverPos, err := Rdb.client.GeoPos(ctx, zsetName, strconv.Itoa(id)).Result()

	if err != nil {
		return nil, err
	}

	if dirverPos[0] == nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return nil, errors.New("data is not found in redis server!!")
	}

	driverLocation := &types.Redis_DriverLocation{}
	driverLocation.Id = id
	driverLocation.Lat = dirverPos[0].Latitude
	driverLocation.Lng = dirverPos[0].Longitude

	return driverLocation, nil
}

func (r Redis) DeleteDriver(id int, zsetName string) error {
	err := Rdb.client.ZRem(ctx, zsetName, strconv.Itoa(id)).Err()
	if err != nil {
		utils.PrintError(DELETE_DATA_FAIL, err)
		return err
	}
	return nil
}

func (r Redis) GetDriverDetail(id int) (*types.Redis_DriverDetail, error) {
	driverInfo := &types.Redis_DriverDetail{}
	err := Rdb.client.HGetAll(ctx, DRIVER_DETAIL+":"+strconv.Itoa(id)).Scan(driverInfo)

	if err != nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return nil, err
	}
	return driverInfo, nil
}

func (r Redis) DeleteDriverDetail(id int) error {
	err := Rdb.client.Unlink(ctx, DRIVER_DETAIL+":"+strconv.Itoa(id)).Err()
	if err != nil {
		utils.PrintError(DELETE_DATA_FAIL, err)
		return err
	}
	return nil
}

func (r Redis) DriverBindOrder(driver *types.Redis_DriverLocation) error {

	_, err := r.GetDriver(driver.Id, AVAILABLE_DRIVER)

	if err != nil {
		if err != redis.Nil {
			utils.PrintError(GET_DATA_FAIL, err)
		}
	}
	err = Rdb.client.ZRem(ctx, AVAILABLE_DRIVER, strconv.Itoa(driver.Id)).Err()

	if err != nil {
		utils.PrintError(DELETE_DATA_FAIL, err)
	}
	fmt.Println(123)

	Rdb.client.GeoAdd(ctx, BUSY_DRIVER, &redis.GeoLocation{
		Name:      strconv.Itoa(driver.Id),
		Longitude: driver.Lng,
		Latitude:  driver.Lat,
	})

	return nil
}

func (r Redis) GetNearDriver(lat float64, lng float64) (*types.Redis_DriverLocation, error) {
	query := redis.GeoSearchQuery{
		Longitude:  lng, // 注意：Redis 順序是 Lng, Lat
		Latitude:   lat,
		Radius:     50,    // 半徑 50
		RadiusUnit: "km",  // 單位 公里
		Sort:       "ASC", // 由近到遠排序
		Count:      1,     // 1位
	}

	searchQuery := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: query,
		WithCoord:      true,
	}

	driverPos, err := Rdb.client.GeoSearchLocation(ctx, AVAILABLE_DRIVER, searchQuery).Result()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(driverPos) == 0 {
		log.Printf(NEAR_DRIVER_NOT_FOUND)
		return nil, errors.New(NEAR_DRIVER_NOT_FOUND)
	}

	id, _ := strconv.Atoi(driverPos[0].Name)

	return &types.Redis_DriverLocation{
		Id:  id,
		Lat: driverPos[0].Latitude,
		Lng: driverPos[0].Longitude,
	}, nil

}

func (r Redis) InsertOrder(order *types.Redis_Order) error {
	_, err := Rdb.client.HSet(ctx, ORDER+":"+strconv.Itoa(order.ID), order).Result()

	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return err
	}
	return nil
}

func (r Redis) GetOrder(id int) (*types.Redis_Order, error) {
	order := &types.Redis_Order{}
	err := Rdb.client.HGetAll(ctx, ORDER+":"+strconv.Itoa(id)).Scan(order)

	if err != nil {
		utils.PrintError(GET_DATA_FAIL, err)
		return nil, err
	}
	return order, nil
}

func (r Redis) InsertSessionKey(sessionId string, cusId int) error {
	sessionData := &types.Redis_Session{
		UserType:   "guest",
		CreateTime: time.Now(),
		CustomerId: cusId,
	}

	_, err := Rdb.client.HSet(ctx, SESSION+":"+sessionId, sessionData).Result()

	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return err
	}

	Rdb.client.Expire(ctx, SESSION+":"+sessionId, 1*time.Hour)

	return nil
}

func (r Redis) GetCustomerId(sessionId string) (int, error) {
	data, err := Rdb.client.HGet(ctx, SESSION+":"+sessionId, "customerId").Result()
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(data)

	if err != nil {
		utils.PrintError(CONVERT_DATA_FAIL, err)
		return 0, err
	}
	return id, nil
}

func (r Redis) CheckSessionKey(sessionId string) (bool, error) {
	res, err := Rdb.client.HExists(ctx, SESSION+":"+sessionId, "customerId").Result()

	if err != nil {
		utils.PrintError(INSERT_DATA_FAIL, err)
		return false, err
	}

	Rdb.client.Expire(ctx, SESSION+":"+sessionId, 1*time.Hour)

	return res, err
}

func (r Redis) UpdateDriverLocation(driver *types.Redis_DriverLocation) error {
	fmt.Println(driver.Status)
	if driver.Status != DRIVER_OFF {
		if driver.Status == BUSY_DRIVER {
			Rdb.client.ZRem(ctx, AVAILABLE_DRIVER, strconv.Itoa(driver.Id))
			r.InsertDriver(driver, driver.Status)
		} else {
			Rdb.client.ZRem(ctx, BUSY_DRIVER, strconv.Itoa(driver.Id))
			r.InsertDriver(driver, driver.Status)
		}
	}

	data, err := json.Marshal(driver)

	if err != nil {
		return err
	}

	err = Rdb.client.Publish(ctx, os.Getenv("DRIVERPUBSUBNAME"), data).Err()

	if err != nil {
		utils.PrintError(PUBLIS_DATA_FAIL, err)
		return err
	}

	return nil
}
