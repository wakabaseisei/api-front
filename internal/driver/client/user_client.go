package client

import (
	"context"
	"fmt"
	"net/http"

	"buf.build/gen/go/wakabaseisei/ms-protobuf/connectrpc/go/ms/user/v1/userv1connect"
	"github.com/wakabaseisei/api-front/internal/domain"
	"github.com/wakabaseisei/api-front/internal/domain/service"
	"github.com/wakabaseisei/api-front/internal/driver/client/converter"
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
	res, serr := s.client.CreateUser(ctx, nil)
	if serr != nil {
		return nil, fmt.Errorf("CreateUser: %w", serr)
	}

	return converter.ConvertUserPbToUser(res.Msg), nil
}
