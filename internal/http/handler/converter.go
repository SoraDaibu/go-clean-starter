package handler

import "github.com/SoraDaibu/go-clean-starter/internal/service/user"

func (r *CreateUserRequest) ToCreateUserInput() *user.CreateUserInput {
	return &user.CreateUserInput{
		Name:     r.Name,
		Email:    string(r.Email),
		Password: r.Password,
	}
}
