package auth

import (
	"context"
	"fmt"

	ssopb "github.com/Nikita-Smirnov-idk/go_microservices_template_project/services/sso/v1/contracts/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type Handler struct {
	ssopb.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssopb.RegisterAuthServer(gRPC, &Handler{auth: auth})
}

func (s *Handler) Login(ctx context.Context, req *ssopb.LoginRequest) (*ssopb.LoginResponse, error) {
	if err := ValidateLoginRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error: %v", err))
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssopb.LoginResponse{
		Token: token,
	}, nil
}

func (s *Handler) Register(ctx context.Context, req *ssopb.RegisterRequest) (*ssopb.RegisterResponse, error) {
	if err := ValidateRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error: %v", err))
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssopb.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *Handler) IsAdmin(ctx context.Context, req *ssopb.IsAdminRequest) (*ssopb.IsAdminResponse, error) {
	if err := ValidateIsAdminRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error: %v", err))
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssopb.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
