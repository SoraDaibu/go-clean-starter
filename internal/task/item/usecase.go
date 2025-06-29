package item

import (
	"context"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
)

type ItemTaskUsecase interface {
	ImportItems(ctx context.Context, sourceDir string, dryRun bool) error
}

type itemTaskUsecase struct {
	Tx       repository.Transaction
	ItemRepo domain.ItemRepository
}

// NewItemTaskUsecase creates a new item task usecase
// Following DIP: depends on domain interface, not concrete implementation
func NewItemTaskUsecase(
	tx repository.Transaction,
	itemRepo domain.ItemRepository,
) ItemTaskUsecase {
	return &itemTaskUsecase{
		Tx:       tx,
		ItemRepo: itemRepo,
	}
}
