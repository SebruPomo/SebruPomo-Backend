package main

import (
	"Backend/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	handler.Register(e.Group("/api"))
	e.Start(":1323")
}
