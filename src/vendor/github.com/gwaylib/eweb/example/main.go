package main

import (
	"strings"

	"github.com/gwaylib/eweb"
	_ "github.com/gwaylib/eweb/example/route"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Register router
func init() {
	e := eweb.Default()
	// TODO: fix to env
	e.Debug = true
	// render
	e.Renderer = eweb.GlobTemplate("./public/**/tpl/*.html")

	// middle ware
	e.Use(middleware.Gzip())

	// filter
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			uri := req.URL.Path
			switch {
			case strings.HasPrefix(uri, "/hacheck"):
				return c.String(200, "1")

			case uri == "/":
				// TODO: redirect to need
				// return c.Redirect(301,"/index")
			}

			// next route
			return next(c)
		}
	})

	// static file
	e.Static("/", "./public")
}

func main() {
	// Start server
	e := eweb.Default()
	e.Logger.Fatal(e.Start(":8081"))
}
