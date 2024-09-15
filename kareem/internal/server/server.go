package server

import (
	"github.com/gofiber/fiber/v2"
	"kareem/internal/storage"
)

type FiberServer struct {
	*fiber.App
	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "kareem",
			AppName:      "kareem",
		}),
		db: database.New(),
	}

	return server
}
