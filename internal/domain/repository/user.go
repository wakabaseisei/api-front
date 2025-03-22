package repository

import (
	"context"

	"github.com/wakabaseisei/api-front/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error)
}
