package item

import (
	"context"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
)

type ItemTaskUsecase interface {
	// ImportItems processes rows on a best-effort basis and returns an error if
	// any row fails validation or persistence. Successful rows remain committed;
	// retrying the whole import may create duplicates. A dry run validates rows
	// without persisting them and also reports validation failures.
	ImportItems(ctx context.Context, sourceDir string, dryRun bool) error
}

type itemTaskUsecase struct {
	Tx       repository.Transaction
	ItemRepo domain.ItemRepository
	Logger   logger.Logger
}

// NewItemTaskUsecase creates a new item task usecase
// Following DIP: depends on domain interface, not concrete implementation
func NewItemTaskUsecase(
	tx repository.Transaction,
	itemRepo domain.ItemRepository,
	l logger.Logger,
) ItemTaskUsecase {
	return &itemTaskUsecase{
		Tx:       tx,
		ItemRepo: itemRepo,
		Logger:   l,
	}
}
