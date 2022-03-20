package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sebrupomo/sebrupomo-backend/db"
	"github.com/sebrupomo/sebrupomo-backend/handler"
)

func main() {
	db.GetConnection()

	e := echo.New()
	e.Validator = &handler.CustomValidator{Validator: validator.New()}
	handler.Register(e.Group("/api"))
	e.Start(":1323")
}
