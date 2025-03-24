package converter

import (
	"fmt"

	"github.com/wakabaseisei/api-front/internal/domain"
)

func ConvertUserToGreetMessage(user *domain.User) string {
	return fmt.Sprintf("Hello, %s!", user.Name)
}
