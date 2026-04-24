package database

const (
	//error status
	INSERT_DATA_FAIL = "insert data fail"
	UPDATE_DATA_FAIL = "update data fail"
	DELETE_DATA_FAIL = "delete data fail"
	GET_DATA_FAIL    = "get data fail"

	//driver status
	DRIVER_DETAIL         = "driver_detail"
	AVAILABLE_DRIVER      = "available_driver"
	BUSY_DRIVER           = "busy_driver"
	DRIVER_OFF            = "off"
	DRIVER_NOT_EXIST      = "driver not exist!!"
	NEAR_DRIVER_NOT_FOUND = "near driver not found!!"
	DRIVER_TASK_FINISH    = "finish"

	//order status
	ORDER_CANCEL      = "cancel"
	ORDER_PROCESS     = "process"
	DRIVER_PICKED_UP  = "picked_up"
	ORDER             = "order"
	ORDER_FINISH      = "complete"
	ORDER_CANCEL_FAIL = "order cancel fail"

	//Session status
	SESSION = "session"

	//pub sub
	PUBLIS_DATA_FAIL = "publis data fail"

	// lock
	GET_LOCK_FAIL = "get lock fail"

	// processing logic
	CONVERT_DATA_FAIL = "convert data fail"
)
