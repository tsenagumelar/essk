package usecase

import (
	"context"

	"github.com/tsenagumelar/essk/services/backend/services/auth-service/repositories"
)

type Usecase struct {
	repository repositories.Repository
}

func New(repository repositories.Repository) Usecase {
	return Usecase{repository: repository}
}

func (u Usecase) Health(ctx context.Context) map[string]any {
	return u.repository.Status(ctx)
}

func (u Usecase) Contract(ctx context.Context) map[string]any {
	status := u.repository.Status(ctx)
	status["grpc_service"] = "essk.auth.v1.AuthService"
	status["schema_owner"] = "auth"
	return status
}
