package wsutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetQueryInt(ctx *fiber.Ctx, key string, def, min, max int) (int, error) {
	valStr := ctx.Query(key)
	if valStr == "" {
		return def, nil
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if val < min || (max > 0 && val > max) {
		return 0, fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("value of '%s' must be in bounds [%d, %d]", key, min, max))
	}

	return val, nil
}

func GetQueryBool(ctx *fiber.Ctx, key string, def bool) (bool, error) {
	v := ctx.Query(key)

	switch strings.ToLower(v) {
	case "":
		return def, nil
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fiber.NewError(fiber.StatusBadRequest, "invalid boolean value")
	}
}
