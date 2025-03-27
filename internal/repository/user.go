package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/wakabaseisei/api-front/internal/domain"
)

type userRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *userRepository {
	return &userRepository{
		conn: conn,
	}
}

func (r *userRepository) Create(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error) {
	// TODO: Replace actual DB data later.
	return &domain.User{
		ID:        cmd.ID,
		Name:      cmd.Name,
		CreatedAt: time.Now(),
	}, nil
}
