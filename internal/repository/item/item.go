package item

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
	"github.com/SoraDaibu/go-clean-starter/internal/repository/common"
	"github.com/SoraDaibu/go-clean-starter/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// itemRepository implements domain.ItemRepository
// Following DIP: depends on abstractions (domain interfaces) not concrete implementations
// Following composition: uses BaseRepository for common functionality
type itemRepository struct {
	*repository.BaseRepository
}

// NewItemRepository creates a new item repository implementation
// Following DIP: returns domain interface, not concrete type
func NewItemRepository(pool *pgxpool.Pool) domain.ItemRepository {
	return &itemRepository{
		BaseRepository: repository.NewBaseRepository(pool),
	}
}

// GetItem implements domain.ItemReader
func (r *itemRepository) GetItem(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	queries := r.GetQueries(ctx)
	item, err := queries.GetItem(ctx, common.UUIDToPgtype(id))
	if err != nil {
		return nil, err
	}

	itemID, err := common.PgtypeToUUID(item.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID: %w", err)
	}

	typeID, err := common.Int32PtrToUint(item.TypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid type_id for item %s: %w", id, err)
	}

	return domain.ItemFromSource(itemID, typeID), nil
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
		itemID, err := common.PgtypeToUUID(item.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid item ID: %w", err)
		}
		typeID, err := common.Int32PtrToUint(item.TypeID)
		if err != nil {
			return nil, fmt.Errorf("invalid type_id for item %s: %w", itemID, err)
		}
		result[i] = domain.ItemFromSource(itemID, typeID)
	}

	return result, nil
}

// CreateItem implements domain.ItemWriter
func (r *itemRepository) CreateItem(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	queries := r.GetQueries(ctx)
	createdItem, err := queries.CreateItem(ctx, sqlc.CreateItemParams{
		ID:     common.UUIDToPgtype(item.ID()),
		TypeID: common.UintToInt32Ptr(item.TypeID()),
	})
	if err != nil {
		return nil, err
	}

	itemID, err := common.PgtypeToUUID(createdItem.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid created item ID: %w", err)
	}

	typeID, err := common.Int32PtrToUint(createdItem.TypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid type_id for created item %s: %w", itemID, err)
	}

	return domain.ItemFromSource(itemID, typeID), nil
}

// UpdateItem implements domain.ItemWriter
func (r *itemRepository) UpdateItem(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	queries := r.GetQueries(ctx)
	updatedItem, err := queries.UpdateItem(ctx, sqlc.UpdateItemParams{
		ID:     common.UUIDToPgtype(item.ID()),
		TypeID: common.UintToInt32Ptr(item.TypeID()),
	})
	if err != nil {
		return nil, err
	}

	itemID, err := common.PgtypeToUUID(updatedItem.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid updated item ID: %w", err)
	}

	typeID, err := common.Int32PtrToUint(updatedItem.TypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid type_id for updated item %s: %w", itemID, err)
	}

	return domain.ItemFromSource(itemID, typeID), nil
}

// DeleteItem implements domain.ItemWriter
func (r *itemRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	queries := r.GetQueries(ctx)
	return queries.DeleteItem(ctx, common.UUIDToPgtype(id))
}
