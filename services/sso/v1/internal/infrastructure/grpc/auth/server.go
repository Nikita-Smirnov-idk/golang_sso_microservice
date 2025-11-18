package auth

import ssopb "github.com/Nikita-Smirnov-idk/go_microservices_template_project/pkg/gen/go/sso/v1"

type serverAPI struct {
	ssopb.UnimplementedAuthServer
}
