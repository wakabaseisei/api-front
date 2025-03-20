package grpc

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"

	apifrontv1 "github.com/wakabaseisei/ms-protobuf/gen/go/ms/apifront/v1"
)

func (s *APIFrontService) Greet(
	ctx context.Context,
	req *connect.Request[apifrontv1.GreetRequest],
) (*connect.Response[apifrontv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&apifrontv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
