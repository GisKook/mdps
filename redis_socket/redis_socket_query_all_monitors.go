package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/base"
	"log"
	"strconv"
	"strings"
)

const (
	MONITOR_KEY_SEP   string = ":"
	MONITOR_VALUE_SEP string = "+"
	MONITOR_TIMESTAMP string = "TIMESTAMP"
)

func (socket *RedisSocket) QueryAllMonitors() []*base.RouterMonitor {
	conn := socket.GetConn()
	defer func() {
		conn.Close()
	}()
	var value interface{}
	var cursor_keys []interface{}
	var cursor string = "0"
	var keys []string
	var e error
	var router_monitors []*base.RouterMonitor = make([]*base.RouterMonitor, 0)

	for {
		value, e = conn.Do("SCAN", cursor)
		base.CheckError(e)
		cursor_keys, e = redis.Values(value, e)
		base.CheckError(e)
		cursor, e = redis.String(cursor_keys[0], nil)
		base.CheckError(e)
		keys, e = redis.Strings(cursor_keys[1], nil)
		keys = base.FilterStringArray(keys, PREFIX_MONITORS)
		base.CheckError(e)
		router_monitors = append(router_monitors, socket.PipelineGetMonitorsValue(keys)...)
		if cursor == "0" {
			return router_monitors
		}
	}
}

func (socket *RedisSocket) PipelineGetMonitorsValue(keys []string) []*base.RouterMonitor {
	if len(keys) != 0 {
		conn := socket.GetConn()
		defer func() {
			conn.Close()
		}()

		router_monitors := make([]*base.RouterMonitor, len(keys))

		var index int = 0
		var key string = ""
		for index, key = range keys {
			conn.Send("HGETALL", key)
			router_id_port := strings.Split(key, MONITOR_KEY_SEP)
			router_id, _ := strconv.Atoi(router_id_port[1])
			serial_port, _ := strconv.Atoi(router_id_port[2])
			router_monitors[index] = &base.RouterMonitor{
				RouterID:   uint32(router_id),
				SerialPort: uint8(serial_port),
			}

		}

		conn.Flush()

		for i := 0; i < index+1; i++ {
			v_redis, err := conn.Receive()

			if err != nil {
				log.Println(err)
				continue
			}

			v, _ := redis.ByteSlices(v_redis, nil)
			socket.PipelineSetMonitorValue(v, router_monitors[i])

		}
		conn.Do("")

		return router_monitors

	}

	return nil

}

func (socket *RedisSocket) PipelineSetMonitorValue(raw [][]byte, router_monitor *base.RouterMonitor) {
	item_count := len(raw)
	var index int = 0
	for i := 0; i < item_count/2; i++ {
		key := string(raw[index])
		if key == MONITOR_TIMESTAMP {
			time_stamp, _ := strconv.ParseInt(string(raw[index+1]), 10, 64)
			router_monitor.TimeStamp = time_stamp
		} else {
			modbus_datatype_datalen := strings.Split(key, MONITOR_VALUE_SEP)
			modbus_addr, _ := strconv.Atoi(modbus_datatype_datalen[0])
			datatype, _ := strconv.Atoi(modbus_datatype_datalen[1])
			datalen, _ := strconv.Atoi(modbus_datatype_datalen[2])
			router_monitor.Monitors = append(router_monitor.Monitors,
				&base.RouterMonitorSingle{
					ModbusAddr: uint32(modbus_addr),
					DataType:   uint8(datatype),
					DataLen:    uint8(datalen),
					Data:       raw[index+1],
				})
		}
		index += 2
	}
}
