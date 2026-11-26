// Package repository handles persistence of items in PostgreSQL.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/example/fullstack-starter/backend/internal/model"
)

// ErrNotFound is returned when the requested item does not exist.
var ErrNotFound = errors.New("item not found")

const itemColumns = "id, name, description, created_at, updated_at"

// ItemRepository provides CRUD access to the items table.
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository creates an ItemRepository.
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// List returns every item, newest first.
func (r *ItemRepository) List(ctx context.Context) ([]model.Item, error) {
	query := `SELECT ` + itemColumns + ` FROM items ORDER BY created_at DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]model.Item, 0)
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate items: %w", err)
	}

	return items, nil
}

// GetByID returns a single item or ErrNotFound.
func (r *ItemRepository) GetByID(ctx context.Context, id int64) (model.Item, error) {
	query := `SELECT ` + itemColumns + ` FROM items WHERE id = $1`

	var item model.Item
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Item{}, ErrNotFound
	}
	if err != nil {
		return model.Item{}, fmt.Errorf("get item: %w", err)
	}

	return item, nil
}

// Create inserts a new item and returns the stored row.
func (r *ItemRepository) Create(ctx context.Context, name, description string) (model.Item, error) {
	query := `
		INSERT INTO items (name, description)
		VALUES ($1, $2)
		RETURNING ` + itemColumns

	var item model.Item
	err := r.db.QueryRowContext(ctx, query, name, description).
		Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return model.Item{}, fmt.Errorf("create item: %w", err)
	}

	return item, nil
}

// Update replaces an item's fields and returns the updated row.
func (r *ItemRepository) Update(ctx context.Context, id int64, name, description string) (model.Item, error) {
	query := `
		UPDATE items
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING ` + itemColumns

	var item model.Item
	err := r.db.QueryRowContext(ctx, query, name, description, id).
		Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Item{}, ErrNotFound
	}
	if err != nil {
		return model.Item{}, fmt.Errorf("update item: %w", err)
	}

	return item, nil
}

// Delete removes an item or returns ErrNotFound.
func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read delete result: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
