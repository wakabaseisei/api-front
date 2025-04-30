package service

import (
	"context"

	"github.com/wakabaseisei/api-front/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error)
}
