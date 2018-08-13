package behavior

import (
	"io"

	"github.com/labstack/gommon/color"
)

type consoleAdapter struct {
	c *color.Color
}

// put a log protocol to log queue
func (a *consoleAdapter) Put(e *Event) {
	a.c.Println(string(e.ToJson()))
}

func (a *consoleAdapter) Close() {
}

func NewConsoleClient(out io.Writer) Client {
	c := new(color.Color)
	c.SetOutput(out)
	return &consoleAdapter{
		c: c,
	}
}
