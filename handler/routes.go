package handler

import "github.com/labstack/echo/v4"

func Register(api *echo.Group) {
	profileGroup := api.Group("/profile")
	profileGroup.POST("", SignUp)
}
