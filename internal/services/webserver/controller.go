package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

// Controller is used to register endpoints for
// a specific section of an API
type Controller interface {

	// Setup sets up the controller using the passed
	// di.Container for dependency injection and
	// registers all endpoints to the passed fiber
	// router.
	Setup(container di.Container, router fiber.Router)
}
