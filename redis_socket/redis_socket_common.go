package redis_socket

import (
	"github.com/giskook/mdps/pb"
	"strconv"
)

const (
	PREFIX_ALTERS    string = "TADATA:"
	SEP_ALTERS       string = "+"
	SEP_ALTERS_VALUE string = ","
	SEP_ALTERS_KEY   string = ":"
)

func (socket *RedisSocket) GenAlterKey(alters *Report.DataCommand) string {
	return PREFIX_ALTERS + strconv.FormatUint(alters.Tid, 10) + SEP_ALTERS_KEY + strconv.FormatUint(uint64(alters.SerialPort), 10)
}
