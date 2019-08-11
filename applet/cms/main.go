package main

import (
	"os"
	"strings"
	"text/template"

	"lserver/applet/cms/route"
	"lserver/module/etc"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/eweb"
	"github.com/gwaylib/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var r = eweb.Default()

// 静态目录
func init() {
	// file
	r.Static("/", "./public/")
}

func main() {
	r.Debug = os.Getenv("GIN_MODE") != "release"

	// view
	if eweb.DebugMode() {
		r.Renderer = &eweb.Template{
			template.Must(template.ParseGlob("./public/tpl/index.html")),
		}
	} else {
		r.Renderer = &eweb.Template{
			template.Must(template.ParseGlob("./public/index.html")),
		}
	}

	// middle ware
	r.Use(middleware.Gzip())

	// 过滤器
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			uri := req.URL.Path

			// 静态页面过虑
			switch {
			case strings.HasPrefix(uri, "/hacheck"):
				return c.String(200, "1")
			case
				// 执行下一个router处理
				strings.HasPrefix(uri, "/favicon"),
				strings.HasPrefix(uri, "/robots.txt"),
				strings.HasPrefix(uri, "/download/"),
				strings.HasPrefix(uri, "/img/"),
				strings.HasPrefix(uri, "/css/"),
				strings.HasPrefix(uri, "/js/"),
				strings.HasPrefix(uri, "/vendor/"),
				strings.HasPrefix(uri, "/bower_components/"),
				strings.HasPrefix(uri, "/tpl/"),
				strings.HasPrefix(uri, "/fonts/"),
				strings.HasPrefix(uri, "/l10n/"),
				strings.HasPrefix(uri, "/api/vcode"),
				strings.HasPrefix(uri, route.AccSigninApiPath),
				strings.HasPrefix(uri, route.AccSignupApiPath),
				strings.HasPrefix(uri, route.AccPwdApiPath):
				return next(c)
			}

			// 校验权限
			uc := route.GetUserCache(c)
			// 没登录
			if uc == nil {
				if strings.ToUpper(req.Method) != "GET" {
					return c.String(302, "登录已过期, 请重登录")
				}
				// 跳转登录页面
				return c.Redirect(302, route.AccSigninViewPath)
			}

			// 已登录的情况处理
			if uri == "/" {
				// 如已登录，根请求跳转主页
				return c.Redirect(302, route.IndexPath)
			}

			// 校验更改权限
			if strings.ToUpper(req.Method) != "GET" && !uc.Priv.CheckSum(uri) {
				log.Warn(errors.New("auth fail").As(uc.UserName, uri))
				return c.String(403, "权限不足")
			}

			// 校验通过
			return next(c)
		}
	})

	port := etc.Etc.String("applet/cms", "listen")
	log.Debugf("start at: %s", port)
	if err := r.Start(port); err != nil {
		log.Exit(1, errors.As(err))
		return
	}
}
