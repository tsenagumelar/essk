package usecase

import (
	"context"

	"github.com/tsenagumelar/essk/services/backend/services/iam-service/repositories"
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
	status["schema_owner"] = "iam"
	status["bounded_context"] = []string{"users", "roles", "permissions"}
	return status
}
