package zmq_server

const (
	HTTP_PLC_ID     string = "plc_id"
	HTTP_PLC_SERIAL string = "serial"

	// restart
	HTTP_RESTART       string = "/plc/restart"
	HTTP_RESTART_DELAY string = "delay"

	// rs232 get status
	HTTP_GET_RS232_STATUS             string = "/plc/get_rs232_status"
	HTTP_GET_RS232_STATUS_SERIAL_PORT string = "serial_port"

	// rs232 control command
	//	HTTP_RS232_CONTROL             string = "/plc/control_rs232"
	//	HTTP_RS232_CONTROL_SERIAL_PORT string = "serial_port"
	//	HTTP_RS232_CONTROL_COMMAND     string = "command"
	//	HTTP_RS232_CONTROL_PORT_TYPE   string = "port_type"
	//	HTTP_RS232_CONTROL_IP          string = "ip"
	//	HTTP_RS232_CONTROL_PORT        string = "port"

	// rs485 get status
	HTTP_GET_RS485_STATUS             string = "/plc/get_rs485_status"
	HTTP_GET_RS485_STATUS_SERIAL_PORT string = "serial_port"

	// rs485 control command
	//	HTTP_RS485_CONTROL             string = "/plc/control_rs485"
	//	HTTP_RS485_CONTROL_SERIAL_PORT string = "serial_port"
	//	HTTP_RS485_CONTROL_COMMAND     string = "command"
	//	HTTP_RS485_CONTROL_PORT_TYPE   string = "port_type"
	//	HTTP_RS485_CONTROL_IP          string = "ip"
	//	HTTP_RS485_CONTROL_PORT        string = "port"

	// modbus get status
	HTTP_MODBUS_STATUS                string = "/plc/modbus"
	HTTP_MODBUS_STATUS_PORT           string = "modbus_port"
	HTTP_MODBUS_STATUS_STANDARD_FRAME string = "standard_frame"

	// data download
	HTTP_DATA_DOWNLOAD             string = "/plc/data_download"
	HTTP_DATA_DOWNLOAD_SERIAL_PORT string = "serial_port"
	HTTP_DATA_DOWNLOAD_MODBUS_ADDR string = "modbus_addr"
	HTTP_DATA_DOWNLOAD_DATA_TYPE   string = "data_type"
	HTTP_DATA_DOWNLOAD_DATA        string = "data"

	// data query
	HTTP_DATA_QUERY             string = "/plc/data_query"
	HTTP_DATA_QUERY_SERIAL_PORT string = "serial_port"
	HTTP_DATA_QUERY_MODBUS_ADDR string = "modbus_addr"
	HTTP_DATA_QUERY_DATA_TYPE   string = "data_type"
	HTTP_DATA_QUERY_DATA        string = "data"

	// batch add monitor
	HTTP_BATCH_ADD_MONITOR                         string = "/plc/batch_add_monitor"
	HTTP_BATCH_ADD_MONITOR_SERIAL_PORT             string = "serial_port"
	HTTP_BATCH_ADD_MONITOR_COUNT                   string = "count"
	HTTP_BATCH_ADD_MONITOR_MONITORS                string = "monitors"
	HTTP_BATCH_ADD_MONITOR_MONITORS_ID             string = "id"
	HTTP_BATCH_ADD_MONITOR_MONITORS_MODBUS_ADDR    string = "modbus_addr"
	HTTP_BATCH_ADD_MONITOR_MONITORS_REGISTER_TYPE  string = "register_type"
	HTTP_BATCH_ADD_MONITOR_MONITORS_REGISTER_COUNT string = "register_count"

	// batch add alter
	HTTP_BATCH_ADD_ALTER                       string = "/plc/batch_add_alter"
	HTTP_BATCH_ADD_ALTER_SERILA_PORT           string = "serial_port"
	HTTP_BATCH_ADD_ALTER_COUNT                 string = "count"
	HTTP_BATCH_ADD_ALTER_ALTERS                string = "alters"
	HTTP_BATCH_ADD_ALTER_ALTERS_ID             string = "id"
	HTTP_BATCH_ADD_ALTER_ALTERS_MODBUS_ADDR    string = "modbus_addr"
	HTTP_BATCH_ADD_ALTER_ALTERS_REGISTER_TYPE  string = "register_type"
	HTTP_BATCH_ADD_ALTER_ALTERS_REGISTER_COUNT string = "register_count"
	HTTP_BATCH_ADD_ALTER_ALTERS_COND           string = "cond"
	HTTP_BATCH_ADD_ALTER_ALTERS_THRESHOLD      string = "threshold"

	//////////////RESPONSE////////////////
	//HTTP_RESPONSE_RESULT               string = "result"
	HTTP_RESPONSE_RESULT_SUCCESS       uint8 = 0
	HTTP_RESPONSE_RESULT_SERVER_FAILED uint8 = 1
	HTTP_RESPONSE_RESULT_PARAMTER_ERR  uint8 = 2
	HTTP_RESPONSE_RESULT_TIMEOUT       uint8 = 3

	//HTTP_RESPONSE_SERIAL_PORT string = "serial_port"
)

var HTTP_RESULT []string = []string{"成功", "失败,路由器反馈失败 或 dps服务器内部错误", "参数错误", "超时,路由器掉线 或 路由器反馈慢"}
