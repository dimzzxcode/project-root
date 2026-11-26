// Package service contains the business logic for items.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/example/fullstack-starter/backend/internal/model"
	"github.com/example/fullstack-starter/backend/internal/repository"
)

const (
	maxNameLength        = 255
	maxDescriptionLength = 2000
)

// ErrNotFound is returned when an item does not exist.
var ErrNotFound = errors.New("item not found")

// ValidationError describes invalid input fields.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// ItemService implements the item use cases.
type ItemService struct {
	repo *repository.ItemRepository
}

// NewItemService creates an ItemService.
func NewItemService(repo *repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// List returns all items.
func (s *ItemService) List(ctx context.Context) ([]model.Item, error) {
	return s.repo.List(ctx)
}

// Get returns a single item.
func (s *ItemService) Get(ctx context.Context, id int64) (model.Item, error) {
	item, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Item{}, ErrNotFound
	}
	if err != nil {
		return model.Item{}, fmt.Errorf("get item: %w", err)
	}
	return item, nil
}

// Create validates the payload and stores a new item.
func (s *ItemService) Create(ctx context.Context, payload model.ItemPayload) (model.Item, error) {
	payload, err := validate(payload)
	if err != nil {
		return model.Item{}, err
	}

	item, err := s.repo.Create(ctx, payload.Name, payload.Description)
	if err != nil {
		return model.Item{}, fmt.Errorf("create item: %w", err)
	}
	return item, nil
}

// Update validates the payload and replaces an existing item.
func (s *ItemService) Update(ctx context.Context, id int64, payload model.ItemPayload) (model.Item, error) {
	payload, err := validate(payload)
	if err != nil {
		return model.Item{}, err
	}

	item, err := s.repo.Update(ctx, id, payload.Name, payload.Description)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Item{}, ErrNotFound
	}
	if err != nil {
		return model.Item{}, fmt.Errorf("update item: %w", err)
	}
	return item, nil
}

// Delete removes an item.
func (s *ItemService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

func validate(payload model.ItemPayload) (model.ItemPayload, error) {
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Description = strings.TrimSpace(payload.Description)

	fields := make(map[string]string)

	switch {
	case payload.Name == "":
		fields["name"] = "name is required"
	case utf8.RuneCountInString(payload.Name) > maxNameLength:
		fields["name"] = fmt.Sprintf("name must be at most %d characters", maxNameLength)
	}

	if utf8.RuneCountInString(payload.Description) > maxDescriptionLength {
		fields["description"] = fmt.Sprintf("description must be at most %d characters", maxDescriptionLength)
	}

	if len(fields) > 0 {
		return model.ItemPayload{}, &ValidationError{Fields: fields}
	}

	return payload, nil
}
