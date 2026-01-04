package builder

import (
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/http/handler/user"
	"github.com/SoraDaibu/go-clean-starter/internal/repository"
	itemRepo "github.com/SoraDaibu/go-clean-starter/internal/repository/item"
	userRepo "github.com/SoraDaibu/go-clean-starter/internal/repository/user"
	userUsecase "github.com/SoraDaibu/go-clean-starter/internal/service/user"
	"github.com/SoraDaibu/go-clean-starter/internal/task/item"
)

// InitializeDependency creates a new Dependency instance with all required dependencies
func InitializeDependency(cfg *config.Config) (*Dependency, error) {
	return Resolve(cfg, NewDependencyNeedsAllTrue())
}

// InitializeUserUsecase creates a new UserUsecase instance
func InitializeUserUsecase(d *Dependency) userUsecase.UserUsecase {
	userRepository := userRepo.NewUserRepository(d.DB)
	return userUsecase.NewUserUsecase(userRepository)
}

// InitializeUserHandler creates a new UserHandler instance
func InitializeUserHandler(d *Dependency) *user.UserHandler {
	userRepository := userRepo.NewUserRepository(d.DB)
	uu := userUsecase.NewUserUsecase(userRepository)
	return user.NewUserHandler(uu)
}

// InitializeItemTaskUsecase creates a new ItemTaskUsecase instance
func InitializeItemTaskUsecase(d *Dependency) item.ItemTaskUsecase {
	transaction := repository.NewTransaction(d.DB)
	itemRepository := itemRepo.NewItemRepository(d.DB)
	return item.NewItemTaskUsecase(transaction, itemRepository)
}
