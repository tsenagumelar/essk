package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db     *pgxpool.Pool
	schema string
}

func New(db *pgxpool.Pool, schema string) Repository {
	return Repository{db: db, schema: schema}
}

func (r Repository) Status(ctx context.Context) map[string]any {
	status := map[string]any{"service": "audit-service", "schema": r.schema, "database": "not_configured"}
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
