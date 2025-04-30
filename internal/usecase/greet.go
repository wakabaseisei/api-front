package usecase

import (
	"context"
	"fmt"

	"github.com/wakabaseisei/api-front/internal/domain"
	"github.com/wakabaseisei/api-front/internal/domain/service"
)

type GreetInteractor interface {
	Invoke(ctx context.Context, cmd *domain.UserCommand) (*domain.User, error)
}

type geetInteractor struct {
	userService service.UserService
}

func NewGreetInteractor(userRepo service.UserService) GreetInteractor {
	return &geetInteractor{
		userService: userRepo,
	}
}

func (i *geetInteractor) Invoke(
	ctx context.Context,
	cmd *domain.UserCommand,
) (
	*domain.User,
	error,
) {
	user, rerr := i.userService.CreateUser(ctx, cmd)
	if rerr != nil {
		return nil, fmt.Errorf("userService Create: %v", rerr)
	}

	return user, nil
}
