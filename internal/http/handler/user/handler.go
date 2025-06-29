package user

import (
	"github.com/SoraDaibu/go-clean-starter/internal/service/user"
)

type UserHandler struct {
	usecase user.UserUsecase
}

func NewUserHandler(
	usecase user.UserUsecase,
) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}
