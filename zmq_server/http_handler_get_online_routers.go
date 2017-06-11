package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/redis_socket"
	"net/http"
)

type GetOnlineRoutersResponse struct {
	Result uint8  `json:"result"`
	Desc   string `json:"desc"`
	Count  uint32 `json:"count"`
}

func EncodingGetOnlineRoutersResponse(get_online_routers_resp *GetOnlineRoutersResponse) string {
	response, _ := json.Marshal(get_online_routers_resp)

	return string(response)
}

func GetOnlineRoutersHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))

		}
	}()
	PrintRequest(r)

	res := &GetOnlineRoutersResponse{
		Result: 0,
		Desc:   "成功",
		Count:  uint32(redis_socket.GetStatusChecker().GetNodeCount()),
	}
	fmt.Fprint(w, EncodingGetOnlineRoutersResponse(res))
}
