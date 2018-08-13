package route

import (
	"module/etc"

	"github.com/gwaylib/eweb"
	"github.com/labstack/echo"
)

// 系统信息
var SysInfo = eweb.H{
	"name":    etc.Etc.String("applet/cms", "name"),
	"version": etc.Etc.String("applet/cms", "version"),
}

const (
	IndexPath = "/app/dashboard"
)

func init() {
	r := eweb.Default()

	// view
	r.GET(IndexPath, Index)
}

// 主页面
func Index(c echo.Context) error {
	return c.Render(200, "index.html", SysInfo)
}
