package behavior

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventTime time.Time
	EventKey  string

	FromIp    string
	ReqMethod string
	ReqBody   string
	RespCode  string
	RespBody  string

	UseTime int64
}

func (b *Event) ToJson() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		panic(*b)
	}
	return data
}

func Parse(src []byte) (*Event, error) {
	b := &Event{}
	if err := json.Unmarshal(src, b); err != nil {
		return nil, err
	}
	return b, nil
}

type Client interface {
	Close()
	Put(*Event)
}
