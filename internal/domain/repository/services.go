package repository

type Services struct {
	userRepository UserRepository
}

func NewServices(userRepository UserRepository) *Services {
	return &Services{
		userRepository: userRepository,
	}
}
