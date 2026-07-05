package usecase

import (
	"context"
	"time"

	"github.com/tsenagumelar/essk/services/backend/services/api-gateway/repositories"
)

type Usecase struct {
	repository repositories.Repository
	upstreams  []string
}

func New(repository repositories.Repository, upstreams []string) Usecase {
	if len(upstreams) == 0 {
		upstreams = []string{
			"127.0.0.1:19110",
			"127.0.0.1:19120",
			"127.0.0.1:19130",
			"127.0.0.1:19140",
			"127.0.0.1:19150",
		}
	}
	return Usecase{repository: repository, upstreams: upstreams}
}

func (u Usecase) Health(ctx context.Context) map[string]any {
	return u.repository.Status(ctx)
}

func (u Usecase) Upstreams(ctx context.Context) map[string]any {
	status := u.repository.Status(ctx)
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	upstreams := make([]map[string]any, 0, len(u.upstreams))
	for _, target := range u.upstreams {
		result := map[string]any{"target": target}
		if err := u.repository.CheckGRPCHealth(checkCtx, target); err != nil {
			result["status"] = "degraded"
			result["error"] = err.Error()
		} else {
			result["status"] = "serving"
		}
		upstreams = append(upstreams, result)
	}
	status["grpc_upstreams"] = upstreams
	return status
}
