package client

import (
	"context"
	"fmt"
	"net/http"

	"buf.build/gen/go/wakabaseisei/ms-protobuf/connectrpc/go/ms/user/v1/userv1connect"
	userv1 "buf.build/gen/go/wakabaseisei/ms-protobuf/protocolbuffers/go/ms/user/v1"
	"connectrpc.com/connect"
	"github.com/wakabaseisei/api-front/internal/domain"
	"github.com/wakabaseisei/api-front/internal/domain/service"
	"github.com/wakabaseisei/api-front/internal/driver/client/converter"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	client userv1connect.UserServiceClient
}

func NewUserService(endpoint string) service.UserService {
	cli := userv1connect.NewUserServiceClient(
		http.DefaultClient,
		endpoint,
	)

	return &userService{
		client: cli,
	}
}

func (s *userService) CreateUser(
	ctx context.Context,
	cmd *domain.UserCommand,
) (
	*domain.User,
	error,
) {
	user := userv1.User{
		UserId:    cmd.ID,
		Name:      cmd.Name,
		CreatedAt: timestamppb.New(cmd.CreatedAt),
	}

	res, serr := s.client.CreateUser(ctx, connect.NewRequest(
		&userv1.CreateUserRequest{
			User: &user,
		},
	))
	if serr != nil {
		return nil, fmt.Errorf("CreateUser: %w", serr)
	}

	return converter.ConvertUserPbToUser(res.Msg), nil
}
