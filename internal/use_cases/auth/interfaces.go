package auth

import (
	"context"

	"github.com/Nikita-Smirnov-idk/go_microservices_template_project/services/sso/v1/internal/domain/models"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int) (models.App, error)
}
