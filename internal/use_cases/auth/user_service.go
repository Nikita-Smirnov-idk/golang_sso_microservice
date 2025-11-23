package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Nikita-Smirnov-idk/go_microservices_template_project/services/sso/v1/internal/infrastructure/repository"
	"github.com/Nikita-Smirnov-idk/go_microservices_template_project/services/sso/v1/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type UserService struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func NewUserService(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *UserService {
	return &UserService{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (s *UserService) Login(ctx context.Context, email string, password string, appID int) (token string, err error) {
	const op = "Auth.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	app, err := s.appProvider.GetApp(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			log.Warn("app not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}
		log.Error("failed to get app", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Info("invalid credentials", "error", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err = jwt.NewToken(user, app, s.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	return token, nil
}

func (s *UserService) Register(ctx context.Context, email string, password string) (userID int64, err error) {
	const op = "userService.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering User")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "error", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			log.Warn("user already exists", "error", err)

			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		log.Error("failed to save new user", "error", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (s *UserService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("user_uid", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := s.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("user not found", "error", err)

			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to check if user is admin", "error", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
