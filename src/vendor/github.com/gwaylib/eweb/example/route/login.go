package routes

import (
	"github.com/gwaylib/eweb"
	"github.com/labstack/echo"
)

// register routes handle
func init() {
	e := eweb.Default()
	e.GET("/login", LoginView)
}

// Login Handler
func LoginView(c echo.Context) error {
	return c.HTML(200, "test/tpl/login")
}
