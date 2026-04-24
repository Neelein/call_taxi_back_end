package api

const (
	AVAILABLE_DRIVER = "available_driver"
	BUSY_DRIVER      = "busy_driver"
	DRIVER_WORKING   = "working"
	DRIVER_OFF       = "off"

	//http error
	INVALID_DATA      = "invalid data format"
	DATA_IS_EMPTY     = "data is empty"
	SERVER_ERROR      = "server error"
	DATA_IS_NOT_EXSIT = "data is not exsit"

	//data processing error
	READ_REQUEST_DATA_FAIL = "read request data fail"
	CONVERT_DATA_FAIL      = "convert data fail"
	BODY_IS_EMPTY          = "body is empty"

	//driver
	INSERT_DRIVER_FAIL          = "insert driver fail"
	INSERT_DRIVER_SUCCESS       = "insert driver success"
	UPDATE_DRIVER_FAIL          = "update driver fail"
	UPDATE_DRIVER_SUCCESS       = "update driver success"
	GET_DRIVER_DATA_FAIL        = "get driver data fail"
	DRIVER_WORKING_FAIL         = "driver working fail"
	DRIVER_WORKING_SUCCESS      = "driver is working"
	UPDATE_DRIVER_LOCATION_FAIL = "update driver location fail"
	NEAR_DRIVER_NOT_FOUND       = "near driver not found!!"
	CALLTAXI_FAIl               = "call taxi fail"
	DRIVER_TASK_FINISH          = "finish"

	//order
	UPDATE_ORDER_SUCCESS = "update order success"
	UPDATE_ORDER_FAIL    = "update order fail"
	GET_ORDER_DATA_FAIL  = "get order data fail"
	ORDER_FINISH         = "order is complete!!"

	// customer
	CUSTOMER_DATA_NOT_EXIST = "customer data not exist"

	//session
	CREATE_SESSION_SUCCESS = "create session success"
	SESSION_KEY_NAME       = "session_id"
	SESSION_ID_NOT_EXIST   = "session id not exist"
)
