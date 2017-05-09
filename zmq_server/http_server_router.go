package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"log"
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
		case key := <-h.HttpRequestDel:
			close(h.HttpRespones[key])
			delete(h.HttpRespones, key)
		case res := <-h.HttpResponseChan:
			chan_resp, ok := h.HttpRespones[res.Key]
			if ok {
				res.Command.SerialNumber = h.GetSerialID(uint16(res.Command.SerialNumber))
				chan_resp <- res.Command
			} else {
				log.Println("nokey")
			}

		}
	}
}

func (h *Http_server) DoResponse(resp *HttpResponsePair) {
	h.HttpResponseChan <- resp
}

func (h *Http_server) SendRequest(key uint64) chan *Report.ControlCommand {
	chan_response := make(chan *Report.ControlCommand)
	h.AddRequest(&HttpRequestPair{
		Key:     key,
		Command: chan_response,
	})

	return chan_response
}
