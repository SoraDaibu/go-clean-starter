package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/SoraDaibu/go-clean-starter/domain"
)

type UserUsecase interface {
	GetUser(ctx context.Context, id uuid.UUID) (*UserOutput, error)
	CreateUser(ctx context.Context, input *CreateUserInput) (*UserOutput, error)
}

type userUsecase struct {
	userRepository domain.UserRepository
}

// NewUserUsecase creates a new user usecase
// Following DIP: depends on domain interface, not concrete implementation
func NewUserUsecase(userRepository domain.UserRepository) UserUsecase {
	return &userUsecase{userRepository: userRepository}
}
