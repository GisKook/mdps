package redis_socket

import (
	"github.com/HuKeping/rbtree"
	"github.com/giskook/mdps/conf"
	"log"
	"sync"
	"time"
)

type Status_Checker struct {
	Rbt_Tid_Time   *rbtree.Rbtree
	Mutex_Tid_Time sync.Mutex

	Rbt_Time_Tid   *rbtree.Rbtree
	Mutex_Time_Tid sync.Mutex
}

type Tid_Time_Status struct {
	Tid      uint64
	RecvTime int64
}

func (x Tid_Time_Status) Less(than rbtree.Item) bool {
	return x.Tid < than.(Tid_Time_Status).Tid
}

type Time_Tid_Status struct {
	RecvTime int64
	Tids     []uint64
}

func (x Time_Tid_Status) Less(than rbtree.Item) bool {
	return x.RecvTime < than.(Time_Tid_Status).RecvTime
}

var G_Status_Checker *Status_Checker

func GetStatusChecker() *Status_Checker {
	if G_Status_Checker == nil {
		G_Status_Checker = &Status_Checker{
			Rbt_Tid_Time: rbtree.New(),
			Rbt_Time_Tid: rbtree.New(),
		}
	}

	return G_Status_Checker
}

func (sc *Status_Checker) Insert(tid uint64, recv_time_stamp int64) {
	sc.Mutex_Tid_Time.Lock()
	sc.Mutex_Time_Tid.Lock()
	defer func() {
		sc.Mutex_Tid_Time.Unlock()
		sc.Mutex_Time_Tid.Unlock()
	}()
	log.Println("insert")
	//1.insert into rbt_tid_time
	//    1.1 if has ->update timevalue
	//    1.2 if not has -> insert
	tid_time := sc.Rbt_Tid_Time.Get(Tid_Time_Status{
		Tid: tid,
	})
	if tid_time != nil {
		sc.Rbt_Tid_Time.Delete(Tid_Time_Status{
			Tid: tid,
		})
	}
	sc.Rbt_Tid_Time.Insert(Tid_Time_Status{
		Tid:      tid,
		RecvTime: recv_time_stamp,
	})
	// 1. if do not have then add
	time_tid := sc.Rbt_Time_Tid.Get(Time_Tid_Status{
		RecvTime: recv_time_stamp,
	})
	if time_tid == nil {
		sc.Rbt_Time_Tid.Insert(Time_Tid_Status{
			RecvTime: recv_time_stamp,
			Tids:     []uint64{tid},
		})
	} else {
		time_tid_status := time_tid.(Time_Tid_Status)
		time_tid_status.Tids = append(time_tid_status.Tids, tid)
		sc.Rbt_Time_Tid.Delete(Time_Tid_Status{
			RecvTime: recv_time_stamp,
		})

		sc.Rbt_Time_Tid.Insert(Time_Tid_Status{
			RecvTime: recv_time_stamp,
			Tids:     time_tid_status.Tids,
		})
	}

}

func (sc *Status_Checker) Del(recv_time_stamp int64) {
	sc.Mutex_Tid_Time.Lock()
	sc.Mutex_Time_Tid.Lock()
	defer func() {
		sc.Mutex_Tid_Time.Unlock()
		sc.Mutex_Time_Tid.Unlock()
	}()

	//1.del the rbt_time_tid
	time_tid_item := sc.Rbt_Time_Tid.Get(Time_Tid_Status{
		RecvTime: recv_time_stamp,
	})
	sc.Rbt_Time_Tid.Delete(Time_Tid_Status{
		RecvTime: recv_time_stamp,
	})
	tids_status := time_tid_item.(Time_Tid_Status)
	for _, _tid := range tids_status.Tids {
		sc.Rbt_Tid_Time.Delete(Tid_Time_Status{
			Tid: _tid,
		})
	}
}

func (sc *Status_Checker) Min() (int64, []uint64) {
	sc.Mutex_Time_Tid.Lock()
	defer func() {
		sc.Mutex_Time_Tid.Unlock()
	}()

	if sc.Rbt_Time_Tid.Len() == 0 {
		return 0, nil
	}

	time_tid_status := sc.Rbt_Time_Tid.Min().(Time_Tid_Status)

	return time_tid_status.RecvTime, time_tid_status.Tids
}

func (sc *Status_Checker) Check() {
	current_time := time.Now().Unix()
	var recv_time int64
	var tids []uint64
	for {
		recv_time, tids = sc.Min()
		if len(tids) > 0 {
			if current_time-recv_time > int64(conf.GetConf().Redis.StatusExpire) {
				log.Println("add off line")
				for _, tid := range tids {
					GetRedisSocket().Terminal_Status_Chan <- &TStatus{
						Tid:    tid,
						Status: 1,
					}
				}

				sc.Del(recv_time)
			} else {
				return
			}
		} else {
			return
		}
	}
}
