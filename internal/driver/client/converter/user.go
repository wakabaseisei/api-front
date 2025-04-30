package converter

import (
	userv1 "buf.build/gen/go/wakabaseisei/ms-protobuf/protocolbuffers/go/ms/user/v1"

	"github.com/wakabaseisei/api-front/internal/domain"
)

func ConvertUserPbToUser(user *userv1.User) *domain.User {
	return &domain.User{
		ID:        user.GetUserId(),
		Name:      user.GetName(),
		CreatedAt: user.GetCreatedAt().AsTime(),
	}
}
