package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsenagumelar/essk/services/backend/internal/platform/service"
	"github.com/tsenagumelar/essk/services/backend/services/iam-service/handler"
	"github.com/tsenagumelar/essk/services/backend/services/iam-service/repositories"
	"github.com/tsenagumelar/essk/services/backend/services/iam-service/usecase"
)

func Register(api fiber.Router, deps service.Dependencies) {
	repo := repositories.New(deps.DB, deps.Config.Service.DatabaseSchema)
	uc := usecase.New(repo)
	h := handler.New(uc)

	group := api.Group("/iam-service")
	group.Get("/health", h.Health)
	group.Get("/contract", h.Contract)
}
