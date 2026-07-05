package main

import (
	"context"
	"log"

	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	"github.com/tsenagumelar/essk/services/backend/services/audit-service/routes"
)

func main() {
	cfg := config.Load("audit-service", 18150, 19150, "audit")
	if err := service.Run(context.Background(), cfg, routes.Register); err != nil {
		log.Fatal(err)
	}
}
