package repository

import (
	"context"
	"time"

	"github.com/wakabaseisei/api-front/internal/domain"
)

type userRepository struct{}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error) {
	// TODO: Replace actual DB data later.
	return &domain.User{
		ID:        cmd.ID,
		Name:      cmd.Name,
		CreatedAt: time.Now(),
	}, nil
}
