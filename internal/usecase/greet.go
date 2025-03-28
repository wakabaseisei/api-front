package usecase

import (
	"context"
	"log"

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
	// if rerr := i.userRepository.Create(ctx, cmd); rerr != nil {
	// 	return nil, fmt.Errorf("userRepository Create: %v", rerr)
	// }

	// user, ferr := i.userRepository.FindByID(ctx, cmd.ID)
	// if ferr != nil {
	// 	return nil, fmt.Errorf("userRepository FindByID: %v", ferr)
	// }

	// return user, nil
	log.Println("Hererere")
	if err := i.userRepository.Ping(ctx); err != nil {
		return nil, err
	}

	return &domain.User{}, nil
}
