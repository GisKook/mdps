package zmq_server

import (
	"github.com/giskook/mdps/pb"
)

type HttpRequestPair struct {
	Key     uint64
	Command chan *Report.ControlCommand
}

type HttpResponsePair struct {
	Key     uint64
	Command *Report.ControlCommand
}

func (h *Http_server) AddRequest(req *HttpRequestPair) {
	h.HttpRequestAdd <- req
}

func (h *Http_server) DelRequest(key uint64) {
	h.HttpRequestDel <- key
}

func (h *Http_server) Run() {
	for {
		select {
		case add := <-h.HttpRequestAdd:
			h.HttpRespones[add.Key] = add.Command
		case del := <-h.HttpRequestDel:
			delete(h.HttpRespones, del.Key)
		case res := <-h.HttpResponseChan:
			chan_resp, ok := h.HttpRespones[res.Key]
			if ok {
				chan_resp <- res.Command
			}

		}
	}
}

func (h *Http_server) DoResponse(resp *HttpResponsePair) {
	h.HttpResponseChan <- resp
}
