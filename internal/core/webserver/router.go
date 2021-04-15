package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type Router interface {
	SetContainer(container di.Container)
	Route(router fiber.Router)
}
