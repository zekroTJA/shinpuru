package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type Controller interface {
	Setup(container di.Container, router fiber.Router)
}
