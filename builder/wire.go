//go:build wireinject
// +build wireinject

package builder

import (
	"database/sql"

	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/http/handler/user"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
	itemRepo "github.com/SoraDaibu/go-clean-starter/internal/repository/item"
	userRepo "github.com/SoraDaibu/go-clean-starter/internal/repository/user"
	userUsecase "github.com/SoraDaibu/go-clean-starter/internal/service/user"
	"github.com/SoraDaibu/go-clean-starter/internal/task/item"
	"github.com/google/wire"
)

//go:generate wire

var (
	// UserUsecaseSet provides all dependencies for UserUsecase
	UserUsecaseSet = wire.NewSet(
		provideDBFromDependency,
		userRepo.NewUserRepository,
		userUsecase.NewUserUsecase,
	)

	// ItemTaskUsecaseSet provides all dependencies for ItemTaskUsecase
	ItemTaskUsecaseSet = wire.NewSet(
		provideDBFromDependency,
		repository.NewTransaction,
		itemRepo.NewItemRepository,
		item.NewItemTaskUsecase,
	)
)

// InitializeDependency creates a new Dependency instance with all required dependencies
func InitializeDependency(cfg *config.Config) (*Dependency, error) {
	return Resolve(cfg, NewDependencyNeedsAllTrue())
}

// InitializeUserUsecase creates a new UserUsecase instance
func InitializeUserUsecase(d *Dependency) userUsecase.UserUsecase {
	wire.Build(UserUsecaseSet)
	return nil
}

// InitializeUserHandler creates a new UserHandler instance
func InitializeUserHandler(d *Dependency) *user.UserHandler {
	wire.Build(UserUsecaseSet, user.NewUserHandler)
	return nil
}

// InitializeItemTaskUsecase creates a new ItemTaskUsecase instance
func InitializeItemTaskUsecase(d *Dependency) item.ItemTaskUsecase {
	wire.Build(ItemTaskUsecaseSet)
	return nil
}

func provideDBFromDependency(d *Dependency) *sql.DB {
	return d.DB
}
