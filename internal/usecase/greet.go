package usecase

import (
	"context"
	"fmt"

	"github.com/wakabaseisei/api-front/internal/domain"
	"github.com/wakabaseisei/api-front/internal/domain/repository"
)

type GreetInteractor interface {
	Invoke(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error)
}

type geetInteractor struct {
	userRepository repository.UserRepository
}

func NewGreetInteractor(userRepo repository.UserRepository) GreetInteractor {
	return &geetInteractor{
		userRepository: userRepo,
	}
}

func (i *geetInteractor) Invoke(
	ctx context.Context,
	cmd *domain.UserCommand,
) (
	*domain.User,
	error,
) {
	user, rerr := i.userRepository.Create(ctx, cmd)
	if rerr != nil {
		return nil, fmt.Errorf("userRepository Create: %v", rerr)
	}

	return user, nil
}
