package main

import (
	"context"
	"log"

	"github.com/tsenagumelar/essk/services/backend/internal/platform/config"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	"github.com/tsenagumelar/essk/services/backend/services/iam-service/routes"
)

func main() {
	cfg := config.Load("iam-service", 18130, 19130, "iam")
	if err := service.Run(context.Background(), cfg, routes.Register); err != nil {
		log.Fatal(err)
	}
}
