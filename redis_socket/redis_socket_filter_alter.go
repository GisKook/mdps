package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/base"
	"github.com/giskook/mdps/pb"
	"strconv"
	"strings"
)

func (socket *RedisSocket) FilterAlters(alters []*Report.DataCommand, alters_redis []*base.RouterAlterRedis) []*base.RouterAlterDB {

}
