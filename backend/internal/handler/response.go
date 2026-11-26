// Package handler exposes the HTTP layer built on Fiber.
package handler

import "github.com/gofiber/fiber/v3"

// APIResponse is the envelope used by every endpoint.
type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *APIError `json:"error,omitempty"`
}

// APIError describes a failure returned by the API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func respond(c fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

func respondError(c fiber.Ctx, status int, code, message string, details any) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
