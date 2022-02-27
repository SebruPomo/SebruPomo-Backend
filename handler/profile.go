package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SignUp(c echo.Context) error {
	return c.JSON(http.StatusCreated, "")
}
