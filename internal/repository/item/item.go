package item

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
	"github.com/SoraDaibu/go-clean-starter/internal/repository/common"
	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
)

// itemRepository implements domain.ItemRepository
// Following DIP: depends on abstractions (domain interfaces) not concrete implementations
// Following composition: uses BaseRepository for common functionality
type itemRepository struct {
	*repository.BaseRepository
}

// NewItemRepository creates a new item repository implementation
// Following DIP: returns domain interface, not concrete type
func NewItemRepository(db *sql.DB) domain.ItemRepository {
	return &itemRepository{
		BaseRepository: repository.NewBaseRepository(db),
	}
}

// GetItem implements domain.ItemReader
func (r *itemRepository) GetItem(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	queries := r.GetQueries(ctx)
	item, err := queries.GetItem(ctx, id)
	if err != nil {
		return nil, err
	}

	typeID, err := common.SqlNullInt32ToUint(item.TypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid type_id for item %s: %w", id, err)
	}

	return domain.ItemFromSource(item.ID, typeID), nil
}

// ListItems implements domain.ItemReader
// Note: The current sqlc query doesn't support limit/offset, so we apply manual pagination
func (r *itemRepository) ListItems(ctx context.Context, limit, offset int) ([]*domain.Item, error) {
	queries := r.GetQueries(ctx)
	items, err := queries.ListItems(ctx)
	if err != nil {
		return nil, err
	}

	// Apply manual pagination since sqlc query doesn't support it
	start := offset
	end := offset + limit
	if start >= len(items) {
		return []*domain.Item{}, nil
	}
	if end > len(items) {
		end = len(items)
	}

	result := make([]*domain.Item, end-start)
	for i, item := range items[start:end] {
		typeID, err := common.SqlNullInt32ToUint(item.TypeID)
		if err != nil {
			return nil, fmt.Errorf("invalid type_id for item %s: %w", item.ID, err)
		}
		result[i] = domain.ItemFromSource(item.ID, typeID)
	}

	return result, nil
}

// CreateItem implements domain.ItemWriter
func (r *itemRepository) CreateItem(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	queries := r.GetQueries(ctx)
	createdItem, err := queries.CreateItem(ctx, sqlc.CreateItemParams{
		ID:     item.ID(),
		TypeID: common.UintToSqlNullInt32(item.TypeID()),
	})
	if err != nil {
		return nil, err
	}

	typeID, err := common.SqlNullInt32ToUint(createdItem.TypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid type_id for created item %s: %w", createdItem.ID, err)
	}

	return domain.ItemFromSource(createdItem.ID, typeID), nil
}

// UpdateItem implements domain.ItemWriter
func (r *itemRepository) UpdateItem(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	queries := r.GetQueries(ctx)
	updatedItem, err := queries.UpdateItem(ctx, sqlc.UpdateItemParams{
		ID:     item.ID(),
		TypeID: common.UintToSqlNullInt32(item.TypeID()),
	})
	if err != nil {
		return nil, err
	}

	typeID, err := common.SqlNullInt32ToUint(updatedItem.TypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid type_id for updated item %s: %w", updatedItem.ID, err)
	}

	return domain.ItemFromSource(updatedItem.ID, typeID), nil
}

// DeleteItem implements domain.ItemWriter
func (r *itemRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	queries := r.GetQueries(ctx)
	return queries.DeleteItem(ctx, id)
}
