package auth

import (
	"context"

	"github.com/Nikita-Smirnov-idk/golang_sso_microservice/internal/domain/models"
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
