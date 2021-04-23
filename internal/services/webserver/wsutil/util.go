package wsutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
)

// GetQueryInt tries to get a value from request query
// and transforms it to an integer value.
//
// If the query value is not provided, def is returened.
//
// If the integer value is smaller than min or larger
// than max (if max is larger than 0), a bounds error
// is returned.
//
// Returned errors are in form of fiber errors with
// appropriate error codes.
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

// GetQueryBool tries to get a value from request query
// and transforms it to an bool value.
//
// If the query value is not provided, def is returened.
//
// Valid string values for <true> are 'true', '1' or
// 'yes. Valid values for <false> are 'false', '0'
// or 'no'.
//
// Returned errors are in form of fiber errors with
// appropriate error codes.
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

// ErrInternalOrNotFound returns a fiber not found
// error when the passed error is a ErrDatabaseNotFound
// error. Otherwise, the passed error is returned
// unchanged.
func ErrInternalOrNotFound(err error) error {
	if database.IsErrDatabaseNotFound(err) {
		return fiber.ErrNotFound
	}
	return err
}
