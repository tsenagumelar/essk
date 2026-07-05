package main

import (
	"context"
	"log"

	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	"github.com/tsenagumelar/essk/services/backend/services/auth-service/routes"
)

func main() {
	cfg := config.Load("auth-service", 18110, 19110, "auth")
	if err := service.Run(context.Background(), cfg, routes.Register); err != nil {
		log.Fatal(err)
	}
}
