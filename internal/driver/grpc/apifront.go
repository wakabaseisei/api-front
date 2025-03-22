package grpc

import "github.com/wakabaseisei/api-front/internal/domain/repository"

type APIFrontService struct {
	services *repository.Services
}

func NewAPIFrontService(services *repository.Services) *APIFrontService {
	return &APIFrontService{
		services: services,
	}
}
