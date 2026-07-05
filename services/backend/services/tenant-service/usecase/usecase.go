package usecase

import (
	"context"

	"github.com/tsenagumelar/essk/services/backend/services/tenant-service/repositories"
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
	status["grpc_service"] = "essk.tenant.v1.TenantService"
	status["schema_owner"] = "tenant"
	return status
}
