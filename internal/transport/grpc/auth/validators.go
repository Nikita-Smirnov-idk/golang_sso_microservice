package auth

import (
	ssopb "github.com/Nikita-Smirnov-idk/go_microservices_template_project/services/sso/v1/contracts/gen/go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func ValidateLoginRequest(req *ssopb.LoginRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(req.GetEmail(), validation.Required, is.Email),
		validation.Field(req.GetPassword(), validation.Required, validation.Length(8, 100)),
		validation.Field(req.GetAppId(), validation.Required, validation.NotIn(0)),
	)
}

func ValidateRegisterRequest(req *ssopb.RegisterRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(req.GetEmail(), validation.Required, is.Email),
		validation.Field(req.GetPassword(), validation.Required, validation.Length(8, 100)),
	)
}

func ValidateIsAdminRequest(req *ssopb.IsAdminRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(req.GetUserId(), validation.Required, validation.NotIn(0)),
	)
}
