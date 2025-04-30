package service

type Services struct {
	UserService UserService
}

func NewServices(userService UserService) *Services {
	return &Services{
		UserService: userService,
	}
}
