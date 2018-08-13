package behavior

import (
	"github.com/gwaylib/beanmsq"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/logger"
)

type lserverAdapter struct {
	p beanmsq.Producer
}

// put a log protocol to log queue
func (a *lserverAdapter) Put(e *Event) {
	if err := a.p.Put(e.ToJson()); err != nil {
		logger.FailLog(errors.As(err, *e))
		return
	}
}

func (a *lserverAdapter) Close() {
	if err := a.p.Close(); err != nil {
		logger.FailLog(err)
	}
}

func NewLServerClient(addr, tube string, pool int) Client {
	return &lserverAdapter{
		p: beanmsq.NewProducer(pool, addr, tube),
	}
}
