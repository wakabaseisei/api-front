package grpc

import (
	"context"
	"log"

	"connectrpc.com/connect"

	apifrontv1 "github.com/wakabaseisei/ms-protobuf/gen/go/ms/apifront/v1"
)

func (s *APIFrontService) Ping(
	ctx context.Context,
	req *connect.Request[apifrontv1.PingRequest],
) (*connect.Response[apifrontv1.PingResponse], error) {
	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&apifrontv1.PingResponse{})
	res.Header().Set("Ping-Version", "v1")
	return res, nil
}
