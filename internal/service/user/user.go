package user

import (
	"context"

	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/google/uuid"
)

func (u *userUsecase) GetUser(ctx context.Context, id uuid.UUID) (*UserOutput, error) {
	user, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return NewUserOutput(user), nil
}

func (u *userUsecase) CreateUser(ctx context.Context, input *CreateUserInput) (*UserOutput, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	user, err := domain.NewUser(input.Name, input.Email, domain.Password(input.Password))
	if err != nil {
		return nil, err
	}

	createdUser, err := u.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return NewUserOutput(createdUser), nil
}
