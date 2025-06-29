package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/SoraDaibu/go-clean-starter/internal/http/base"
	"github.com/SoraDaibu/go-clean-starter/internal/http/handler"
)

func (u *UserHandler) GetUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return base.HandleError(c, err)
	}

	user, err := u.usecase.GetUser(c.Request().Context(), userID)
	if err != nil {
		return base.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, handler.UserResponse{
		Id:   user.ID,
		Name: user.Name,
	})
}

func (u *UserHandler) CreateUser(c echo.Context) error {
	var req handler.CreateUserRequest
	if err := base.Bind(c, &req); err != nil {
		return err
	}

	user, err := u.usecase.CreateUser(c.Request().Context(), req.ToCreateUserInput())
	if err != nil {
		return base.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, handler.UserResponse{
		Id:   user.ID,
		Name: user.Name,
	})
}
