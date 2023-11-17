package handler

import (
	"admin-api/handler/dto"
	"admin-api/repository"
	"admin-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	AdminHandler interface {
		CreateAdminUser(c echo.Context) error
		SignIn(c echo.Context) error
		RefreshToken(c echo.Context) error
		GetAllAdminEmail(c echo.Context) error
	}
	adminHandler struct {
		ucAdmin   usecase.AdminUsecase
		repoAdmin repository.AdminUserRepository
	}
)

func NewAdminHandler() AdminHandler {
	return &adminHandler{
		ucAdmin:   usecase.NewAdminUsecase(),
		repoAdmin: repository.NewAdminUserRepository(),
	}
}

func (a *adminHandler) CreateAdminUser(c echo.Context) error {
	var req dto.CreateAdminUserRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := a.ucAdmin.CreateAdminUser(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Create admin user successfully!!",
		"data":    req.Email,
	})
}

func (a *adminHandler) SignIn(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := a.ucAdmin.SignIn(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (a *adminHandler) RefreshToken(c echo.Context) error {
	var req dto.RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := a.ucAdmin.RefreshToken(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (a *adminHandler) GetAllAdminEmail(c echo.Context) error {
	emails, err := a.repoAdmin.GetAllAdminEmail()
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &dto.GetAllAdminEmailResponse{
		Emails: emails,
	})
}
