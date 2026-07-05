package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Repository struct {
	db     *pgxpool.Pool
	schema string
}

func New(db *pgxpool.Pool, schema string) Repository {
	return Repository{db: db, schema: schema}
}

func (r Repository) Status(ctx context.Context) map[string]any {
	status := map[string]any{"service": "api-gateway", "schema": r.schema, "database": "not_configured"}
	if r.db == nil {
		return status
	}
	if err := r.db.Ping(ctx); err != nil {
		status["database"] = "degraded"
		status["error"] = err.Error()
		return status
	}
	status["database"] = "ok"
	return status
}

func (r Repository) CheckGRPCHealth(ctx context.Context, target string) error {
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	_, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	return err
}
