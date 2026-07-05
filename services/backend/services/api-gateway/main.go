package main

import (
	"context"
	"log"

	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	"github.com/tsenagumelar/essk/services/backend/services/api-gateway/routes"
)

func main() {
	cfg := config.Load("api-gateway", 18080, 19100, "gateway")
	if err := service.Run(context.Background(), cfg, routes.Register); err != nil {
		log.Fatal(err)
	}
}
