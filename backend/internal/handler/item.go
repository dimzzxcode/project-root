package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/example/fullstack-starter/backend/internal/model"
	"github.com/example/fullstack-starter/backend/internal/service"
	"github.com/gofiber/fiber/v3"
)

// ItemHandler exposes the item CRUD endpoints.
type ItemHandler struct {
	service *service.ItemService
	logger  *slog.Logger
}

// NewItemHandler creates an ItemHandler.
func NewItemHandler(itemService *service.ItemService, logger *slog.Logger) *ItemHandler {
	return &ItemHandler{service: itemService, logger: logger}
}

// Register mounts the item routes on the given router.
func (h *ItemHandler) Register(router fiber.Router) {
	items := router.Group("/items")
	items.Get("/", h.List)
	items.Post("/", h.Create)
	items.Get("/:id", h.Get)
	items.Put("/:id", h.Update)
	items.Delete("/:id", h.Delete)
}

// List handles GET /api/items.
func (h *ItemHandler) List(c fiber.Ctx) error {
	items, err := h.service.List(c.Context())
	if err != nil {
		return h.internalError(c, err)
	}
	return respond(c, fiber.StatusOK, items)
}

// Get handles GET /api/items/:id.
func (h *ItemHandler) Get(c fiber.Ctx) error {
	id, ok := h.parseID(c)
	if !ok {
		return nil
	}

	item, err := h.service.Get(c.Context(), id)
	if err != nil {
		return h.handleServiceError(c, err)
	}
	return respond(c, fiber.StatusOK, item)
}

// Create handles POST /api/items.
func (h *ItemHandler) Create(c fiber.Ctx) error {
	payload, ok := h.parsePayload(c)
	if !ok {
		return nil
	}

	item, err := h.service.Create(c.Context(), payload)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	c.Set(fiber.HeaderLocation, fmt.Sprintf("/api/items/%d", item.ID))
	return respond(c, fiber.StatusCreated, item)
}

// Update handles PUT /api/items/:id.
func (h *ItemHandler) Update(c fiber.Ctx) error {
	id, ok := h.parseID(c)
	if !ok {
		return nil
	}

	payload, ok := h.parsePayload(c)
	if !ok {
		return nil
	}

	item, err := h.service.Update(c.Context(), id, payload)
	if err != nil {
		return h.handleServiceError(c, err)
	}
	return respond(c, fiber.StatusOK, item)
}

// Delete handles DELETE /api/items/:id.
func (h *ItemHandler) Delete(c fiber.Ctx) error {
	id, ok := h.parseID(c)
	if !ok {
		return nil
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleServiceError(c, err)
	}
	return respond(c, fiber.StatusOK, nil)
}

func (h *ItemHandler) parseID(c fiber.Ctx) (int64, bool) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		_ = respondError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "id must be a positive integer", map[string]string{
			"id": "must be a positive integer",
		})
		return 0, false
	}
	return id, true
}

func (h *ItemHandler) parsePayload(c fiber.Ctx) (model.ItemPayload, bool) {
	var payload model.ItemPayload
	if err := c.Bind().Body(&payload); err != nil {
		_ = respondError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "request body must be valid JSON", nil)
		return model.ItemPayload{}, false
	}
	return payload, true
}

func (h *ItemHandler) handleServiceError(c fiber.Ctx, err error) error {
	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		return respondError(c, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", "validation failed", validationErr.Fields)
	}
	if errors.Is(err, service.ErrNotFound) {
		return respondError(c, fiber.StatusNotFound, "NOT_FOUND", "item not found", nil)
	}
	return h.internalError(c, err)
}

func (h *ItemHandler) internalError(c fiber.Ctx, err error) error {
	h.logger.Error("request failed",
		"method", c.Method(),
		"path", c.Path(),
		"error", err,
	)
	return respondError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil)
}
