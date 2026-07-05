package main

import (
	"context"
	"log"

	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	"github.com/tsenagumelar/essk/services/backend/services/tenant-service/routes"
)

func main() {
	cfg := config.Load("tenant-service", 18120, 19120, "tenant")
	if err := service.Run(context.Background(), cfg, routes.Register); err != nil {
		log.Fatal(err)
	}
}
