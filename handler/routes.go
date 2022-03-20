package handler

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/jwt"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func Register(api *echo.Group) {
	jwt.Setup()

	profileGroup := api.Group("/profile")
	profileGroup.POST("", SignUp)
	profileGroup.POST("/login", Login)

	profileRestrictedGroup := api.Group("/profile")
	profileRestrictedGroup.Use(jwt.JwtMiddlware)
	profileRestrictedGroup.GET("/whoAmI", WhoAmI)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

type Context struct {
	echo.Context
}

func (c *Context) BindValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}
	return nil
}

func (c *Context) BindValidateParams(i interface{}) error {
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, i); err != nil {
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}
	return nil
}
