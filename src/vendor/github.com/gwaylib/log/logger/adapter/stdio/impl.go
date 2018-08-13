package stdio

import (
	"io"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/proto"
	"github.com/labstack/gommon/color"
)

type Adapter struct {
	c *color.Color
}

// put a log protocol to log queue
func (a *Adapter) Put(p *proto.Proto) {
	if p == nil {
		panic("argument is nil")
	}
	for _, val := range p.Data {
		date := val.Date.Format("2006-01-02 15:04:05.000")
		a.c.Printf("%s %-5s [%s] %s",
			date,
			val.Level.ColorString(),
			color.Cyan(val.Logger),
			string(val.Msg),
		)
	}
}

func (a *Adapter) Close() {
}

func New(out io.Writer) logger.Adapter {
	c := new(color.Color)
	c.SetOutput(out)
	return &Adapter{
		c: c,
	}
}
