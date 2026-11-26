package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

// NewErrorHandler converts Fiber and unexpected errors into the API response
// envelope. Internal error details are logged but never sent to clients.
func NewErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := "internal server error"

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
			switch status {
			case fiber.StatusNotFound:
				code = "NOT_FOUND"
				message = "route not found"
			case fiber.StatusMethodNotAllowed:
				code = "METHOD_NOT_ALLOWED"
				message = "method not allowed"
			default:
				code = "HTTP_ERROR"
				message = fiberErr.Message
			}
		}

		if status >= fiber.StatusInternalServerError {
			logger.Error("unhandled request error",
				"method", c.Method(),
				"path", c.Path(),
				"error", err,
			)
			code = "INTERNAL_ERROR"
			message = "internal server error"
		}

		return respondError(c, status, code, message, nil)
	}
}
