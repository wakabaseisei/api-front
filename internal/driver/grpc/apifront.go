package grpc

import "github.com/wakabaseisei/api-front/internal/domain/service"

type APIFrontService struct {
	services *service.Services
}

func NewAPIFrontService(services *service.Services) *APIFrontService {
	return &APIFrontService{
		services: services,
	}
}
