package httptry

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gwaylib/errors"
)

const (
	DefaultSleepTime = 3 * 1e9
	DefaultTryTimes  = 3
	DefaultTimeOut   = 20 * 1e9 // 服务器之间不应出现过久的响应
)

var (
	DefaultClient = NewClient(nil, DefaultTryTimes, DefaultSleepTime, DefaultTimeOut)

	SSLClient = NewClient(
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS10}, // SSLv3存在安全问题，最低使用1.0
				Proxy:           http.ProxyFromEnvironment,
			},
		},
		DefaultTryTimes, DefaultSleepTime, DefaultTimeOut,
	)
)

var (
	ErrReqTimeout = errors.New("request timeout")
)

func Do(req *http.Request) (*http.Response, error) {
	return DefaultClient.Do(req)
}

func Get(url string) (r *http.Response, err error) {
	return DefaultClient.Get(url)
}

func Post(url, bodyType string, body io.Reader) (r *http.Response, err error) {
	return DefaultClient.Post(url, bodyType, body)
}

func PostForm(url string, val url.Values) (r *http.Response, err error) {
	return DefaultClient.PostForm(url, val)
}

type Client struct {
	client    *http.Client
	tryTimes  int
	sleepTime time.Duration
	timeout   time.Duration
}

func NewClient(client *http.Client, tryTimes int, sleepTime, timeout time.Duration) *Client {
	if client == nil {
		client = &http.Client{
			Timeout: timeout,
		}
	}
	return &Client{
		client:    client,
		tryTimes:  tryTimes,
		sleepTime: sleepTime,
		timeout:   timeout,
	}
}

func (c *Client) Do(req *http.Request) (r *http.Response, err error) {
	for i := c.tryTimes; i > 0; i-- {
		r, err = c.doTimeout(req)
		if err == nil && r.StatusCode == http.StatusOK {
			return r, nil
		}

		// if it has err and the body is exist,
		// do close the body.
		if r != nil && r.Body != nil {
			r.Body.Close()
		}

		// in the last time,not need to sleep
		// if sleeptime is zero,not need to sleep
		if i > 1 && c.sleepTime > 0 {
			time.Sleep(c.sleepTime)
		}
	}
	return r, errors.As(err)
}

func (c *Client) Get(url string) (r *http.Response, err error) {
	for i := c.tryTimes; i > 0; i-- {
		r, err = c.getTimeout(url)
		if err == nil && r.StatusCode == http.StatusOK {
			return r, nil
		}

		// if it has err and the body is exist,
		// do close the body.
		if r != nil && r.Body != nil {
			r.Body.Close()
		}

		// in the last time,not need to sleep
		// if sleeptime is zero,not need to sleep
		if i > 1 && c.sleepTime > 0 {
			time.Sleep(c.sleepTime)
		}
	}
	return r, err
}

func (c *Client) PostForm(url string, val url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(val.Encode()))
}

func (c *Client) Post(url, bodyType string, body io.Reader) (r *http.Response, err error) {
	oriData, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	oriBody := strings.NewReader(string(oriData))
	for i := 0; i < c.tryTimes; i++ {
		oriBody.Seek(0, 0)
		r, err = c.postTimeout(url, bodyType, oriBody)
		if err == nil {
			return r, nil
		}

		// if it has err and the body is exist,
		// do close the body.
		if r != nil && r.Body != nil {
			r.Body.Close()
		}

		// in the last time,not need to sleep
		// if sleeptime is zero,not need to sleep
		if i < c.tryTimes-1 && c.sleepTime > 0 {
			time.Sleep(c.sleepTime)
		}
	}
	return
}

func (c *Client) doTimeout(req *http.Request) (r *http.Response, err error) {
	wait := make(chan bool)
	go func() {
		defer close(wait)
		r, err = c.client.Do(req)
		wait <- true
	}()
	select {
	case <-wait:
	case <-time.After(c.timeout):
		err = ErrReqTimeout.As("Do")
	}
	return r, errors.As(err)
}

func (c *Client) getTimeout(url string) (r *http.Response, err error) {
	wait := make(chan bool)
	go func() {
		defer close(wait)
		r, err = c.client.Get(url)
		wait <- true
	}()
	select {
	case <-wait:
	case <-time.After(c.timeout):
		err = ErrReqTimeout.As(url)
	}
	return r, errors.As(err)
}

func (c *Client) postTimeout(url, bodyType string, body io.Reader) (r *http.Response, err error) {
	wait := make(chan bool)
	go func() {
		defer close(wait)
		r, err = c.client.Post(url, bodyType, body)
		wait <- true
	}()
	select {
	case <-wait:
	case <-time.After(c.timeout):
		err = ErrReqTimeout.As(url, bodyType)
	}
	return r, errors.As(err)
}
