package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/wakabaseisei/api-front/internal/domain"
	"github.com/wakabaseisei/api-front/internal/driver/grpc/converter"
	"github.com/wakabaseisei/api-front/internal/usecase"
	apifrontv1 "github.com/wakabaseisei/ms-protobuf/gen/go/ms/apifront/v1"
)

func (s *APIFrontService) Greet(
	ctx context.Context,
	req *connect.Request[apifrontv1.GreetRequest],
) (*connect.Response[apifrontv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())

	cmd := &domain.UserCommand{
		ID:        uuid.NewString(),
		Name:      req.Msg.GetName(),
		CreatedAt: time.Now(),
	}
	user, uerr := usecase.NewGreetInteractor(s.services.UserRepository).Invoke(ctx, cmd)
	if uerr != nil {
		return nil, fmt.Errorf("usecase.GreetInteractor.Invoke(): %v", uerr)
	}

	res := connect.NewResponse(&apifrontv1.GreetResponse{
		Greeting: converter.ConvertUserToGreetMessage(user),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
