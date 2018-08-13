package eweb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	stdLog "log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
)

var (
	// Global instance
	defaultE     *Eweb
	defaultELock = sync.Mutex{}
)

type Template struct {
	*template.Template
}

func NewTemplate(tpl *template.Template) *Template {
	return &Template{tpl}
}
func GlobTemplate(filePath string) *Template {
	return &Template{template.Must(template.ParseGlob(filePath))}
}

// Implements Renderer interface
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}
func DebugMode() bool {
	return defaultE.Debug
}

// Struct for rendering map
type H map[string]interface{}

type Eweb struct {
	*echo.Echo
}

// Using global instance to manager router packages
func Default() *Eweb {
	defaultELock.Lock()
	defer defaultELock.Unlock()
	if defaultE == nil {
		defaultE = &Eweb{
			Echo: echo.New(),
		}

		// fix log to stderr
		//	defaultE.Logger.SetOutput(os.Stderr)
		//	defaultE.Logger.SetPrefix("eweb")
		defaultE.Server.Handler = defaultE
		defaultE.TLSServer.Handler = defaultE

		// monitor middleware
		defaultE.Pre(defaultE.FilterHandle)

		// defaultE.Use(middleware.Recover())
	}
	return defaultE
}

func (e *Eweb) colorForStatus(code string) string {
	switch {
	case code >= "200" && code < "300":
		return color.Green(code)
	case code >= "300" && code < "400":
		return color.White(code)
	case code >= "400" && code < "500":
		return color.Yellow(code)
	default:
		return color.Red(code)
	}
}

func (e *Eweb) colorForMethod(method string) string {
	switch method {
	case "GET":
		return color.Blue(method)
	case "POST":
		return color.Cyan(method)
	case "PUT":
		return color.Yellow(method)
	case "DELETE":
		return color.Red(method)
	case "PATCH":
		return color.Green(method)
	case "HEAD":
		return color.Magenta(method)
	case "OPTIONS":
		return color.White(method)
	default:
		return color.Reset(method)
	}
}

func (e *Eweb) FilterHandle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			resp := c.Response()
			req := c.Request()

			// 回传数据
			req.Header.Add("resp-size", strconv.FormatInt(resp.Size, 10))
			req.Header.Add("resp-status", strconv.Itoa(resp.Status))
		}()

		if err := next(c); err != nil {
			// deal default error
			var (
				code = http.StatusInternalServerError
				msg  string
			)

			if he, ok := err.(*echo.HTTPError); ok {
				code = he.Code
				msg = fmt.Sprintf("%v", he.Message)
			} else if e.Debug {
				msg = err.Error()
			} else {
				msg = http.StatusText(code)
			}
			// Send response
			if !c.Response().Committed {
				if c.Request().Method == echo.HEAD {
					return c.NoContent(code)
				}
				return c.String(code, msg)
			}
			return err
		}
		return nil
	}
}

func (e *Eweb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// request performance
	defer func(start time.Time) {
		pErr := recover()

		if pErr != nil {
			r.Header.Set("resp-status", "500")
		}

		req := r
		n := r.Header.Get("resp-status")
		if !e.Debug && n < "400" {
			return
		}

		ip := ReadIp(r)
		stop := time.Now()
		color.Printf(
			"[eweb] %s | %s | %-4s | USE:%12s | FROM:%15s | URI:%s \n",
			start.Format("2006-01-02 15:04:05"),
			e.colorForStatus(n), req.Method,
			// fmt.Sprintf("%.3fs", float64(stop.Sub(start)/1e6)/1000),
			stop.Sub(start).String(),
			ip, req.RequestURI,
		)

		if pErr != nil {
			fmt.Println("panic: %+v\n", pErr)
			panic(pErr)
		}
	}(time.Now())

	e.Echo.ServeHTTP(w, r)
}

// Start starts an HTTP server.
func (e *Eweb) Start(address string) error {
	e.Server.Addr = address
	return e.StartServer(e.Server)
}

// StartTLS starts an HTTPS server.
func (e *Eweb) StartTLS(address string, certFile, keyFile string) (err error) {
	if certFile == "" || keyFile == "" {
		return errors.New("invalid tls configuration")
	}
	c := new(tls.Config)
	c.Certificates = make([]tls.Certificate, 1)
	c.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	return e.startTLS(address, c)
}

// StartAutoTLS starts an HTTPS server using certificates automatically installed from https://letsencrypt.org.
func (e *Eweb) StartAutoTLS(address string) error {
	go http.ListenAndServe(":http", e.AutoTLSManager.HTTPHandler(nil))

	c := new(tls.Config)
	c.GetCertificate = e.AutoTLSManager.GetCertificate
	return e.startTLS(address, c)
}

func (e Eweb) StartTLSConfig(addr string, c *tls.Config) error {
	return e.startTLS(addr, c)
}

func (e *Eweb) startTLS(addr string, c *tls.Config) error {
	s := e.TLSServer
	s.Addr = addr
	s.TLSConfig = c

	if !e.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return e.StartServer(e.TLSServer)
}

// StartServer starts a custom http server.
func (e *Eweb) StartServer(s *http.Server) (err error) {
	s.ErrorLog = stdLog.New(os.Stderr, "", stdLog.LstdFlags)
	s.Handler = e

	kl, err := newListener(s.Addr)
	if err != nil {
		return err
	}

	if s.TLSConfig != nil {
		// https
		return s.Serve(tls.NewListener(kl, s.TLSConfig))
	}
	return s.Serve(kl)
}
