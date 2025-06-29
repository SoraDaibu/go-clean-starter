package user

import (
	"github.com/SoraDaibu/go-clean-starter/domain"
	"github.com/google/uuid"
)

type UserOutput struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUserOutput(user *domain.User) *UserOutput {
	return &UserOutput{
		ID:   user.ID(),
		Name: user.Name(),
	}
}
