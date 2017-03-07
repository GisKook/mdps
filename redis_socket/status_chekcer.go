package redis_socket

import (
	"github.com/HuKeping/rbtree"
	"github.com/giskook/mdps/conf"
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
	// 1. update tid_time rbtree
	// 2. update time_tid rbtree

	sc.Mutex_Tid_Time.Lock()
	sc.Mutex_Time_Tid.Lock()
	defer func() {
		sc.Mutex_Tid_Time.Unlock()
		sc.Mutex_Time_Tid.Unlock()
	}()
	//1.insert into rbt_tid_time
	//    1.1 if has ->update timevalue
	//    1.2 if not has -> insert
	tid_time := sc.Rbt_Tid_Time.Get(Tid_Time_Status{
		Tid: tid,
	})

	store_time := recv_time_stamp
	if tid_time != nil {
		store_time = tid_time.(Tid_Time_Status).RecvTime
		sc.Rbt_Tid_Time.Delete(tid_time.(Tid_Time_Status))
	}
	sc.Rbt_Tid_Time.Insert(Tid_Time_Status{
		Tid:      tid,
		RecvTime: recv_time_stamp,
	})

	// sync time_tid rbtree
	time_tid := sc.Rbt_Time_Tid.Get(Time_Tid_Status{
		RecvTime: store_time,
	})
	if time_tid == nil {
		sc.Rbt_Time_Tid.Insert(Time_Tid_Status{
			RecvTime: recv_time_stamp,
			Tids:     []uint64{tid},
		})
	} else {
		time_tid_status := time_tid.(Time_Tid_Status)
		for i, _tid := range time_tid_status.Tids {
			if _tid == tid {
				time_tid_status.Tids[i] = time_tid_status.Tids[len(time_tid_status.Tids)-1]
				time_tid_status.Tids = time_tid_status.Tids[:len(time_tid_status.Tids)-1]

				//				if i == 0 {
				//					time_tid_status.Tids = time_tid_status.Tids[1:]
				//				} else if i == len(time_tid_status.Tids) {
				//					time_tid_status.Tids = time_tid_status.Tids[:len(time_tid_status.Tids)-1]
				//				} else {
				//					time_tid_status.Tids[i] = time_tid_status.Tids[len(time_tid_status.Tids)-1]
				//					time_tid_status.Tids = time_tid_status.Tids[:len(time_tid_status.Tids)-1]
				//			}
			}
		}

		sc.Rbt_Time_Tid.Delete(Time_Tid_Status{
			RecvTime: store_time,
		})
		if len(time_tid_status.Tids) > 0 {
			sc.Rbt_Time_Tid.Insert(Time_Tid_Status{
				RecvTime: store_time,
				Tids:     time_tid_status.Tids,
			})
		}

		_time_tid := sc.Rbt_Time_Tid.Get(Time_Tid_Status{
			RecvTime: recv_time_stamp,
		})
		if _time_tid == nil {
			sc.Rbt_Time_Tid.Insert(Time_Tid_Status{
				RecvTime: recv_time_stamp,
				Tids:     []uint64{tid},
			})
		} else {
			new_time_tid_status := _time_tid.(Time_Tid_Status)
			for _, _tid := range new_time_tid_status.Tids {
				if _tid == tid {
					return
				}
			}
			new_time_tid_status.Tids = append(new_time_tid_status.Tids, tid)
			sc.Rbt_Time_Tid.Delete(Time_Tid_Status{
				RecvTime: recv_time_stamp,
			})

			sc.Rbt_Time_Tid.Insert(Time_Tid_Status{
				RecvTime: recv_time_stamp,
				Tids:     new_time_tid_status.Tids,
			})
		}

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
	time_tid_item := sc.Rbt_Time_Tid.Delete(Time_Tid_Status{
		RecvTime: recv_time_stamp,
	})
	if time_tid_item != nil {
		tids_status := time_tid_item.(Time_Tid_Status)
		for _, _tid := range tids_status.Tids {
			sc.Rbt_Tid_Time.Delete(Tid_Time_Status{
				Tid: _tid,
			})
		}
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
